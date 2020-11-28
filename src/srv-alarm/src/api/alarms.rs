use std::convert::Infallible;

use warp::Filter;

use crate::models::Context;
use sunrise_common::alarm::NextAlarm;

/// Combination of all alarm related filters
pub fn filters(
    ctx: Context,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    next(ctx)
}

/// GET /alarms/next
fn next(ctx: Context) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("alarms" / "next")
        .and(warp::get())
        .and(with_state(ctx))
        .and_then(get_next_alarm)
}

async fn get_next_alarm(ctx: Context) -> Result<impl warp::Reply, Infallible> {
    if let Some(next_alarm) = ctx.get_next_alarm_ring() {
        Ok(warp::reply::json(&next_alarm))
    } else {
        Ok(warp::reply::json(&NextAlarm::default()))
    }
}

fn with_state(
    ctx: Context,
) -> impl Filter<Extract = (Context,), Error = std::convert::Infallible> + Clone {
    warp::any().map(move || ctx.clone())
}
