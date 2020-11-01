use crate::database::Db;
use crate::models::{Alarm, Error, ERROR_ALARM_EXISTS, ERROR_ALARM_NOT_FOUND};
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
    let mut alarms = vec![];
    db.read(|db| {
        alarms = db.alarms.values().cloned().collect();
    })
    .expect("Error while listing alarms");
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
    log::info!("create_alarm: {:?}", alarm);

    // Insert alarm in database
    let mut found: bool = false;
    db.write(|db| {
        found = db.alarms.contains_key(&alarm.id);
        if !found {
            db.alarms.insert(alarm.id, alarm.clone());
        }
    })
    .expect("Error while creating alarm");
    db.save().expect("Failed to save database");

    // Alarm with id already exists
    if found {
        let error = Error {
            code: ERROR_ALARM_EXISTS,
        };
        return Ok(warp::reply::with_status(
            warp::reply::json(&error),
            StatusCode::CONFLICT,
        ));
    }

    // Create alarm success
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
    // Set alarm id
    let mut alarm: Alarm = alarm.clone();
    alarm.id = id;

    // Update alarm in database
    let mut found: bool = false;
    db.write(|db| {
        found = db.alarms.contains_key(&alarm.id);
        if found {
            db.alarms.insert(alarm.id, alarm.clone());
        }
    })
    .expect("Error while updating alarm");
    db.save().expect("Failed to save database");

    // Alarm with id not found
    if !found {
        let error = Error {
            code: ERROR_ALARM_NOT_FOUND,
        };
        return Ok(warp::reply::with_status(
            warp::reply::json(&error),
            StatusCode::NOT_FOUND,
        ));
    }

    // Update alarm success
    let json_reply = warp::reply::json(&alarm);
    return Ok(warp::reply::with_status(json_reply, StatusCode::OK));
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

    // Delete alarm from database
    let mut found: bool = false;
    db.write(|db| {
        found = db.alarms.remove(&id).is_some();
    })
    .expect("Error while deleting alarm");
    db.save().expect("Failed to save database");

    // Return result
    if found {
        // Delete alarm success
        Ok(StatusCode::NO_CONTENT)
    } else {
        // Alarm not found
        log::warn!("    -> alarm id not found!");
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
