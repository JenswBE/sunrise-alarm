use std::convert::Infallible;

use warp::Filter;

use crate::player::{Command, Remote};

/// Combination of all alarm related filters
pub fn filters(
    remote: Remote,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    start(remote.clone()).or(stop(remote.clone()))
}

/// PUT /music
fn start(
    remote: Remote,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("music")
        .and(warp::put())
        .and(with_remote(remote))
        .and_then(start_music)
}

async fn start_music(remote: Remote) -> Result<impl warp::Reply, Infallible> {
    remote.send(Command::Start).unwrap();
    Ok(warp::reply())
}

/// DELETE /music
fn stop(
    remote: Remote,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("music")
        .and(warp::delete())
        .and(with_remote(remote))
        .and_then(stop_music)
}

async fn stop_music(remote: Remote) -> Result<impl warp::Reply, Infallible> {
    remote.send(Command::Stop).unwrap();
    Ok(warp::reply())
}

fn with_remote(
    remote: Remote,
) -> impl Filter<Extract = (Remote,), Error = std::convert::Infallible> + Clone {
    warp::any().map(move || remote.clone())
}
