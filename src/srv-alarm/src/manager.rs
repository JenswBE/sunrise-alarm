use chrono::Local;
use futures_util::FutureExt;
use tokio::select;
use tokio::sync::mpsc;
use tokio::time::{self, Duration, Instant};

use crate::models::{Config, State, Status};
use crate::time::update_next_alarms;
use sunrise_common::alarm::{Alarm, NextAction};

pub type Radio = mpsc::UnboundedSender<Action>;

#[derive(Debug, Clone, Eq, PartialEq)]
pub enum Action {
    UpdateSchedule,
    ButtonPressed,
    ButtonLongPressed,
}

pub fn start(state: State, config: Config) -> Radio {
    let (tx, mut rx) = mpsc::unbounded_channel::<Action>();

    let tx2 = tx.clone();
    tokio::spawn(async move {
        let fallback = Duration::from_secs(1);
        let duration = duration_until_next_action(state.clone()).unwrap_or(fallback);
        let mut delay_next_action = time::delay_for(duration).fuse();
        log::debug!("Initial delay set to {:?}s", duration.as_secs());

        loop {
            select! {
                action = rx.recv() => {
                    if action.is_none() {
                        log::error!("All manager radio senders closed. Exiting manager ...");
                        break;
                    }

                    let duration = handle_action(action.unwrap(), state.clone()).await;
                    let delay = calculate_delay(duration);
                    delay_next_action = time::delay_until(delay).fuse();
                }

                _ = &mut delay_next_action => {
                    handle_next_action(state.clone(), config.clone(), tx2.clone()).await;
                    // Delay will be set by handle_action (through UpdateSchedule)
                }
            }
        }
    });

    return tx;
}

/// Calculate the duration until the next action
fn duration_until_next_action(state: State) -> Result<Duration, String> {
    let state = state.lock().unwrap();
    let alarm = state
        .next_alarm_action
        .as_ref()
        .ok_or("No next alarm action")?;
    let dt = alarm
        .next_action_datetime
        .expect("Next alarm action should have a datetime");
    dt.signed_duration_since(Local::now())
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

async fn handle_action(action: Action, state: State) -> Option<Duration> {
    log::debug!("Handle action: {:?}", action);
    match action {
        Action::UpdateSchedule => handle_update_schedule(state.clone()),
        _ => None,
    }
}

fn handle_update_schedule(state: State) -> Option<Duration> {
    {
        // Lock state
        let state = state.lock().unwrap();

        // Skip update if not idle
        // Will be updated after alarm is stopped
        if state.status != Status::Idle {
            return None;
        }

        // Skip reschedule if no next alarm
        if state.next_alarm_action.is_none() {
            return None;
        }
    }

    // Calculate duration
    Some(duration_until_next_action(state).unwrap_or(Duration::from_secs(1)))
}

async fn handle_next_action(state: State, config: Config, radio: Radio) {
    // Fetch next alarm from state
    let next_alarm = {
        let state = state.lock().unwrap();
        state.next_alarm_action.clone()
    };

    // Check if alarm is ready
    let ready = next_alarm
        .as_ref()
        .and_then(|a| a.next_action_datetime)
        .and_then(|d| Some(d <= Local::now()));

    // Alarm is not set or not ready
    if ready.is_none() || Some(false) == ready {
        log::debug!("Handle next action: None");
        update_next_alarms(state, &config.alarm, radio);
        return;
    }

    // Alarm is ready => Fetch alarm from next_alarm
    let next_alarm = next_alarm.unwrap();
    let alarm = {
        let state = state.lock().unwrap();
        state
            .alarms
            .iter()
            .find(|&a| a.id == next_alarm.id)
            .unwrap()
            .clone()
    };

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
