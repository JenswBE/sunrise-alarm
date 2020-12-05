use std::convert::Infallible;

use warp::Filter;

use crate::player::{Command, Remote};

/// Combination of all alarm related filters
pub fn filters(
    remote: Remote,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    post_music(remote.clone())
        .or(delete_music(remote.clone()))
        .or(post_volume_increase(remote))
}

/// POST /music
fn post_music(
    remote: Remote,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("music")
        .and(warp::post())
        .and(with_remote(remote))
        .and_then(start_music)
}

async fn start_music(remote: Remote) -> Result<impl warp::Reply, Infallible> {
    remote.send(Command::Start).unwrap();
    Ok(warp::reply())
}

/// DELETE /music
fn delete_music(
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

/// POST /volume/increase
fn post_volume_increase(
    remote: Remote,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("music")
        .and(warp::post())
        .and(with_remote(remote))
        .and_then(increase_volume)
}

async fn increase_volume(remote: Remote) -> Result<impl warp::Reply, Infallible> {
    remote.send(Command::IncreaseVolume).unwrap();
    Ok(warp::reply())
}

// ==============================================
// =                   HELPERS                  =
// ==============================================

fn with_remote(
    remote: Remote,
) -> impl Filter<Extract = (Remote,), Error = std::convert::Infallible> + Clone {
    warp::any().map(move || remote.clone())
}
