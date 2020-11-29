use chrono::Local;
use futures_util::FutureExt;
use tokio::select;
use tokio::sync::mpsc;
use tokio::time::{self, Duration, Instant};

use crate::models::{Context, Status};
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
        let fallback = Duration::from_secs(1);
        let duration = duration_until_next_action(ctx.clone()).unwrap_or(fallback);
        let mut delay_next_action = time::delay_for(duration).fuse();
        log::debug!("Initial delay set to {:?}s", duration.as_secs());

        loop {
            select! {
                action = rx.recv() => {
                    if action.is_none() {
                        log::error!("All manager radio senders closed. Exiting manager ...");
                        break;
                    }

                    let duration = handle_action(ctx.clone(), action.unwrap()).await;
                    let delay = calculate_delay(duration);
                    delay_next_action = time::delay_until(delay).fuse();
                }

                _ = &mut delay_next_action => {
                    handle_next_action(ctx.clone()).await;
                    // Delay will be set by handle_action (through UpdateSchedule)
                }
            }
        }
    });
}

/// Calculate the duration until the next action
fn duration_until_next_action(ctx: Context) -> Result<Duration, String> {
    let alarm = ctx.get_next_alarms_action().ok_or("No next alarm action")?;
    alarm
        .next_action_datetime
        .signed_duration_since(Local::now())
        .to_std()
        .map_err(|_| "Next action already passed".to_string())
}

/// Calculate next Instant until we have to sleep
fn calculate_delay(duration: Option<Duration>) -> Instant {
    if let Some(duration) = duration {
        log::debug!("Reset delay to {:?}s", duration.as_secs());
        return Instant::now() + duration;
    } else {
        log::debug!("No next action. Force check in 15 minutes");
        return Instant::now() + Duration::from_secs(60 * 15);
    }
}

async fn handle_action(ctx: Context, action: Action) -> Option<Duration> {
    log::debug!("Handle action: {:?}", action);
    match action {
        Action::UpdateSchedule => handle_update_schedule(ctx),
        _ => None,
    }
}

fn handle_update_schedule(ctx: Context) -> Option<Duration> {
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

async fn handle_next_action(ctx: Context) {
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
        NextAction::Ring => handle_next_action_ring(alarm).await,
        NextAction::Skip => handle_next_action_skip(alarm).await,
        _ => (),
    }
}

async fn handle_next_action_ring(_alarm: Alarm) {}

async fn handle_next_action_skip(mut alarm: Alarm) {
    alarm.skip_next = false;
    let url = format!("http://localhost:8001/alarms/{}", alarm.id);
    reqwest::Client::new()
        .put(&url)
        .json(&alarm)
        .send()
        .await
        .unwrap();
}
