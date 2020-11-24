use std::convert::Infallible;

use warp::Filter;

use crate::models::State;
use sunrise_common::alarm::NextAlarm;

/// Combination of all alarm related filters
pub fn filters(
    state: State,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    next(state.clone())
}

/// GET /alarms/next
fn next(state: State) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("alarms" / "next")
        .and(warp::get())
        .and(with_state(state))
        .and_then(get_next_alarm)
}

async fn get_next_alarm(state: State) -> Result<impl warp::Reply, Infallible> {
    let locked_state = state.lock().unwrap();
    if let Some(next_alarm) = &locked_state.next_alarm {
        Ok(warp::reply::json(&next_alarm))
    } else {
        Ok(warp::reply::json(&NextAlarm::default()))
    }
}

fn with_state(
    state: State,
) -> impl Filter<Extract = (State,), Error = std::convert::Infallible> + Clone {
    warp::any().map(move || state.clone())
}
