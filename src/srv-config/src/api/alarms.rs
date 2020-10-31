use crate::models::{Alarm, Db, Error, ERROR_ALARM_EXISTS, ERROR_ALARM_NOT_FOUND};
use std::convert::Infallible;
use uuid::Uuid;
use warp::http::StatusCode;
use warp::Filter;

/// Combination of all alarm related filters
pub fn filters(db: Db) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    list(db.clone())
        .or(create(db.clone()))
        .or(update(db.clone()))
        .or(delete(db))
}

/// GET /alarms
fn list(db: Db) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("alarms")
        .and(warp::get())
        .and(with_db(db))
        .and_then(list_alarms)
}

async fn list_alarms(db: Db) -> Result<impl warp::Reply, Infallible> {
    // Just return a JSON array of alarms, applying the limit and offset.
    let alarms = db.lock().await;
    let alarms: Vec<Alarm> = alarms.clone().into_iter().collect();
    Ok(warp::reply::json(&alarms))
}

/// POST /alarms with JSON body
fn create(db: Db) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("alarms")
        .and(warp::post())
        .and(json_body())
        .and(with_db(db))
        .and_then(create_alarm)
}

async fn create_alarm(alarm: Alarm, db: Db) -> Result<impl warp::Reply, Infallible> {
    log::debug!("create_alarm: {:?}", alarm);

    let mut vec = db.lock().await;

    for item in vec.iter() {
        if item.id == alarm.id {
            log::debug!("    -> id already exists: {}", alarm.id);
            // Alarm with id already exists
            let error = Error {
                code: ERROR_ALARM_EXISTS,
            };
            return Ok(warp::reply::with_status(
                warp::reply::json(&error),
                StatusCode::CONFLICT,
            ));
        }
    }

    // No existing Alarm with id, so insert and return `201 Created`.
    vec.push(alarm.clone());

    let json_reply = warp::reply::json(&alarm);
    Ok(warp::reply::with_status(json_reply, StatusCode::CREATED))
}

/// PUT /alarms/:id with JSON body
fn update(db: Db) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("alarms" / Uuid)
        .and(warp::put())
        .and(json_body())
        .and(with_db(db))
        .and_then(update_alarm)
}

async fn update_alarm(id: Uuid, alarm: Alarm, db: Db) -> Result<impl warp::Reply, Infallible> {
    // Lock database
    let mut vec = db.lock().await;

    // Set id to alarm
    let mut alarm: Alarm = alarm.clone();
    alarm.id = id;

    // Update alarm
    for item in vec.iter_mut() {
        if item.id == id {
            *item = alarm.clone();
            let json_reply = warp::reply::json(&alarm);
            return Ok(warp::reply::with_status(json_reply, StatusCode::OK));
        }
    }

    // Alarm not found
    let error = Error {
        code: ERROR_ALARM_NOT_FOUND,
    };
    let json_reply = warp::reply::json(&error);
    return Ok(warp::reply::with_status(json_reply, StatusCode::NOT_FOUND));
}

/// DELETE /alarms/:id
fn delete(db: Db) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("alarms" / Uuid)
        .and(warp::delete())
        .and(with_db(db))
        .and_then(delete_alarm)
}

pub async fn delete_alarm(id: Uuid, db: Db) -> Result<impl warp::Reply, Infallible> {
    log::debug!("delete_alarm: id={}", id);

    let mut vec = db.lock().await;

    let len = vec.len();
    vec.retain(|alarm| {
        // Retain all Alarms that aren't this id...
        // In other words, remove all that *are* this id...
        alarm.id != id
    });

    // If the vec is smaller, we found and deleted a Alarm!
    let deleted = vec.len() != len;

    if deleted {
        Ok(StatusCode::NO_CONTENT)
    } else {
        log::debug!("    -> alarm id not found!");
        Ok(StatusCode::NOT_FOUND)
    }
}

fn with_db(db: Db) -> impl Filter<Extract = (Db,), Error = std::convert::Infallible> + Clone {
    warp::any().map(move || db.clone())
}

fn json_body() -> impl Filter<Extract = (Alarm,), Error = warp::Rejection> + Clone {
    // When accepting a body, we want a JSON body
    // (and to reject huge payloads)...
    warp::body::content_length_limit(1024 * 16).and(warp::body::json())
}
