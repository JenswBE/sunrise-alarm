use futures_util::FutureExt;
use tokio::select;
use tokio::sync::mpsc;
use tokio::time::{self, Duration, Instant};

use crate::models::{Context, Status};

pub type Ringer = mpsc::UnboundedSender<Action>;

#[derive(Debug, Clone, Eq, PartialEq)]
pub enum Action {
    Start,
    Stop,
}

pub fn start(ctx: Context) -> Ringer {
    let (tx, mut rx) = mpsc::unbounded_channel::<Action>();
    tokio::spawn(async move {
        let mut delay_next_step = time::delay_until(Instant::now()).fuse();

        loop {
            select! {
                action = rx.recv() => {
                    if action.is_none() {
                        log::error!("All ringer radio senders closed. Exiting ringer ...");
                        break;
                    }

                    let duration = handle_action(&ctx, action.unwrap()).await;
                    if let Some(duration) = duration {
                        log::debug!("Reset delay to {:?}s", duration.as_secs());
                        delay_next_step = time::delay_until(Instant::now() + duration).fuse();
                    }
                }

                _ = &mut delay_next_step => {
                    // Skip if alarm is idle
                    if ctx.get_status() != Status::Idle {
                        handle_next_step(&ctx).await;
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
        Action::Start => handle_start(ctx),
        Action::Stop => handle_stop(ctx),
    }
}

fn handle_start(_ctx: &Context) -> Option<Duration> {
    Some(Duration::from_secs(60))
}

fn handle_stop(_ctx: &Context) -> Option<Duration> {
    None
}

async fn handle_next_step(_ctx: &Context) {}
