use chrono::Local;
use futures_util::FutureExt;
use tokio::select;
use tokio::sync::mpsc;
use tokio::time::{self, Duration, Instant};

use crate::http::{self, Leds};
use crate::models::{Context, Status};
use crate::ringer::{self, Action as RingerAction, Ringer};
use crate::time::update_next_alarms;
use sunrise_common::alarm::{Alarm, NextAction};

pub type Radio = mpsc::UnboundedSender<Action>;
pub type RadioReceiver = mpsc::UnboundedReceiver<Action>;

#[derive(Debug, Clone, Eq, PartialEq)]
pub enum Action {
    UpdateSchedule,
    ButtonPressed,
    ButtonLongPressed,
}

pub fn create_radios() -> (Radio, RadioReceiver) {
    mpsc::unbounded_channel::<Action>()
}

pub fn start(ctx: Context, mut rx: RadioReceiver) {
    tokio::spawn(async move {
        // Setup initial delay
        let fallback = Duration::from_secs(1);
        let duration = duration_until_next_action(&ctx).unwrap_or(fallback);
        let mut delay_next_action = time::delay_for(duration).fuse();
        log::debug!("Initial delay set to {:?}s", duration.as_secs());

        // Setup ringer
        let ringer = ringer::start(ctx.clone());

        // Start manager loop
        loop {
            select! {
                action = rx.recv() => {
                    if action.is_none() {
                        log::error!("All manager radio senders closed. Exiting manager ...");
                        break;
                    }

                    let duration = handle_action(&ctx, &ringer, action.unwrap()).await;
                    if let Some(duration) = duration {
                        log::debug!("Reset delay to {:?}s", duration.as_secs());
                        delay_next_action = time::delay_until(Instant::now() + duration).fuse();
                    }
                }

                _ = &mut delay_next_action => {
                    handle_next_action(&ctx, &ringer).await;
                    // Delay will be set by handle_action (through UpdateSchedule)
                }
            }
        }
    });
}

/// Calculate the duration until the next action
fn duration_until_next_action(ctx: &Context) -> Result<Duration, String> {
    let alarm = ctx.get_next_alarms_action().ok_or("No next alarm action")?;
    alarm
        .next_action_datetime
        .signed_duration_since(Local::now())
        .to_std()
        .map_err(|_| "Next action already passed".to_string())
}

async fn handle_action(ctx: &Context, ringer: &Ringer, action: Action) -> Option<Duration> {
    log::debug!("Handle action: {:?}", action);
    match action {
        Action::UpdateSchedule => handle_update_schedule(ctx),
        Action::ButtonPressed => handle_button_pressed(ctx, ringer).await,
        Action::ButtonLongPressed => handle_button_long_pressed(ctx, ringer).await,
    }
}

fn handle_update_schedule(ctx: &Context) -> Option<Duration> {
    // Skip update if not idle
    // Will be updated after alarm is stopped
    if ctx.get_status() != Status::Idle {
        return None;
    }

    // Skip reschedule if no next alarm
    if ctx.get_next_alarms_action().is_none() {
        return None;
    }

    // Calculate duration
    Some(duration_until_next_action(ctx).unwrap_or(Duration::from_secs(1)))
}

async fn handle_button_pressed(ctx: &Context, ringer: &Ringer) -> Option<Duration> {
    if ctx.get_status() == Status::Idle {
        // Handle night light
        let leds = http::get_leds(ctx).await.unwrap_or_default();
        if leds.is_off() {
            let req = Leds::night_light();
            http::set_leds(ctx, &req).await.ok();
        } else {
            http::set_leds_off(ctx).await.ok();
        }
    } else {
        // Stop ringer
        ctx.set_status(Status::Idle);
        ringer
            .send(RingerAction::Stop)
            .map_err(|e| log::error!("Failed to stop ringer: {}", e))
            .ok();
    }
    return None;
}

async fn handle_button_long_pressed(ctx: &Context, ringer: &Ringer) -> Option<Duration> {
    if ctx.get_status() == Status::Idle {
        // Handle night light
        let leds = http::get_leds(ctx).await.unwrap_or_default();
        if leds.is_off() {
            let req = Leds::night_light_dark();
            http::set_leds(ctx, &req).await.ok();
        } else {
            http::set_leds_off(ctx).await.ok();
        }
    } else {
        // Stop ringer
        ctx.set_status(Status::Idle);
        ringer
            .send(RingerAction::Stop)
            .map_err(|e| log::error!("Failed to stop ringer: {}", e))
            .ok();
    }
    return None;
}

async fn handle_next_action(ctx: &Context, ringer: &Ringer) {
    // Fetch next alarm from state
    let next_alarm = ctx.get_next_alarms_action();

    // Check if alarm is ready
    let ready = next_alarm
        .as_ref()
        .map(|a| a.next_action_datetime)
        .and_then(|d| Some(d <= Local::now()));

    // Alarm is not set or not ready
    if ready.is_none() || Some(false) == ready {
        log::debug!("Handle next action: None");
        update_next_alarms(ctx);
        return;
    }

    // Alarm is ready => Fetch alarm from next_alarm
    let next_alarm = next_alarm.unwrap();
    let alarm = ctx.get_alarm(next_alarm.id).unwrap();

    // Perform action
    log::debug!("Handle next action: {:?}", next_alarm.next_action);
    match next_alarm.next_action {
        NextAction::Ring => handle_next_action_ring(ctx, ringer, alarm).await,
        NextAction::Skip => handle_next_action_skip(ctx, alarm).await,
        _ => (),
    }
}

async fn handle_next_action_ring(ctx: &Context, ringer: &Ringer, alarm: Alarm) {
    log::debug!("Starting alarm {} ({})", alarm.id, alarm.name);
    ctx.set_status(Status::Ring);
    ringer
        .send(RingerAction::Start)
        .map_err(|e| log::error!("Failed to start ringer: {}", e))
        .ok();
}

async fn handle_next_action_skip(ctx: &Context, mut alarm: Alarm) {
    alarm.skip_next = false;
    http::update_alarm(ctx, alarm)
        .await
        .map_err(|e| log::error!("Failed to update alarm: {}", e))
        .ok();
}
