use futures_util::FutureExt;
use tokio::sync::mpsc;
use tokio::time::{self, Duration, Instant};
use tokio::{join, select};

use crate::http;
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
        let mut delay_next_step = time::delay_until(Instant::now()).fuse();
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
                        log::debug!("Reset delay to {:?}s", duration.as_secs());
                        delay_next_step = time::delay_until(Instant::now() + duration).fuse();
                    }
                }

                _ = &mut delay_next_step => {
                    // Skip if alarm is idle
                    minute += 1;
                    if ctx.get_status() != Status::Idle {
                        log::debug!("Handle next step");
                        handle_next_step(&ctx, minute).await;
                        log::debug!("Reset delay to {}s", TICK_DURATION.as_secs());
                        delay_next_step = time::delay_until(Instant::now() + TICK_DURATION).fuse();
                    } else {
                        log::debug!("Next step skipped since status is idle")
                    }
                }
            }
        }
    });
    return tx;
}

async fn handle_action(ctx: &Context, action: Action) -> Option<Duration> {
    log::debug!("Handle action: {:?}", action);
    match action {
        Action::Start => handle_start(ctx).await,
        Action::Stop => handle_stop(ctx).await,
    }
}

async fn handle_start(ctx: &Context) -> Option<Duration> {
    // Start light
    log::debug!("Starting sunrise");
    http::start_sunrise(ctx).await.ok();
    Some(TICK_DURATION)
}

async fn handle_stop(ctx: &Context) -> Option<Duration> {
    // Stop light and sound
    log::debug!("Stopping sunrise");
    log::debug!("Stopping alarm music");
    let (leds, music) = join!(http::stop_sunrise(ctx), http::stop_music(ctx));
    leds.ok();
    music.ok();

    // Don't set next step
    None
}

async fn handle_next_step(ctx: &Context, minute: u8) {
    let sound_delay = ctx.config.alarm.light_duration - ctx.config.alarm.sound_duration;
    let delay_minutes = sound_delay.num_minutes();
    if minute as i64 == delay_minutes {
        log::debug!("Starting alarm music");
        http::start_music(ctx).await.ok();
    } else if minute as i64 > delay_minutes {
        log::debug!("Increasing alarm volume");
        http::increase_music_volume(ctx).await.ok();
    }
}
