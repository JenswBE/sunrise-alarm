use futures_util::FutureExt;
use tokio::sync::mpsc;
use tokio::time::{self, Duration, Instant};
use tokio::{join, select};

use crate::http;
use crate::manager::Action as MgrAction;
use crate::models::{Context, Status};

pub type Ringer = mpsc::UnboundedSender<Action>;

#[derive(Debug, Clone, Eq, PartialEq)]
pub enum Action {
    Start,
    Stop,
}

const TICK_DURATION: Duration = Duration::from_secs(60);

pub fn start(ctx: Context) -> Ringer {
    let (tx, mut rx) = mpsc::unbounded_channel::<Action>();
    tokio::spawn(async move {
        let mut delay_next_step = Box::pin(time::sleep_until(Instant::now()).fuse());
        let mut minute = 0;

        loop {
            select! {
                action = rx.recv() => {
                    if action.is_none() {
                        log::error!("All ringer radio senders closed. Exiting ringer ...");
                        break;
                    }

                    let duration = handle_action(&ctx, action.unwrap()).await;
                    if let Some(duration) = duration {
                        minute = 0;
                        log::info!("Reset delay to {:?}s", duration.as_secs());
                        delay_next_step = Box::pin(time::sleep_until(Instant::now() + duration).fuse());
                    }
                }

                _ = &mut delay_next_step => {
                    // Skip if alarm is idle
                    minute += 1;
                    if ctx.get_status() != Status::Idle {
                        log::info!("Handle next step");
                        handle_next_step(&ctx, minute).await;
                        log::info!("Reset delay to {}s", TICK_DURATION.as_secs());
                        delay_next_step = Box::pin(time::sleep_until(Instant::now() + TICK_DURATION).fuse());
                    } else {
                        log::info!("Next step skipped since status is idle")
                    }
                }
            }
        }
    });
    return tx;
}

async fn handle_action(ctx: &Context, action: Action) -> Option<Duration> {
    log::info!("Handle action: {:?}", action);
    match action {
        Action::Start => handle_start(ctx).await,
        Action::Stop => handle_stop(ctx).await,
    }
}

async fn handle_start(ctx: &Context) -> Option<Duration> {
    // Start light
    log::info!("Starting sunrise");
    http::start_sunrise(ctx).await.ok();
    Some(TICK_DURATION)
}

async fn handle_stop(ctx: &Context) -> Option<Duration> {
    // Stop light and sound
    log::info!("Stopping sunrise");
    log::info!("Stopping alarm music");
    log::info!("Stopping buzzer");
    let (leds, music, buzzer) = join!(
        http::stop_sunrise(ctx),
        http::stop_music(ctx),
        http::stop_buzzer(ctx)
    );
    leds.ok();
    music.ok();
    buzzer.ok();

    // Don't set next step
    None
}

async fn handle_next_step(ctx: &Context, minute: u8) {
    // General
    let minute = minute as i64;
    log::info!("Handle next alarm step at minute: {}", minute);

    // Check for abort
    let light_duration = ctx.config.alarm.light_duration.num_minutes();
    let abort_delay = light_duration + 10;
    if minute > abort_delay {
        log::warn!("Alarm reached abort limit. Requesting manager to abort alarm.");
        ctx.radio
            .send(MgrAction::AbortAlarm)
            .map_err(|e| log::error!("Failed to send action AbortAlarm to manager: {}", e))
            .ok();
        return;
    }

    // Check for buzzer
    if minute >= light_duration {
        log::info!("Starting buzzer");
        http::start_buzzer(ctx).await.ok();
    }

    // Check for music
    let sound_delay = ctx.config.alarm.light_duration - ctx.config.alarm.sound_duration;
    if minute == sound_delay.num_minutes() {
        log::info!("Starting alarm music");
        http::start_music(ctx).await.ok();
    } else if minute > sound_delay.num_minutes() {
        log::info!("Increasing alarm volume");
        http::increase_music_volume(ctx).await.ok();
    }
}
