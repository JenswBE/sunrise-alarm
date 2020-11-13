use std::convert::Infallible;

use rumqttc::AsyncClient;
use uuid::Uuid;
use warp::http::StatusCode;
use warp::Filter;

use crate::database::Db;
use crate::models::{Error, ERROR_ALARM_EXISTS, ERROR_ALARM_NOT_FOUND};
use crate::mqtt;
use common::general::Alarm;

/// Combination of all alarm related filters
pub fn filters(
    db: Db,
    mqtt_client: AsyncClient,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    list(db.clone())
        .or(create(db.clone(), mqtt_client.clone()))
        .or(update(db.clone(), mqtt_client.clone()))
        .or(delete(db, mqtt_client))
}

/// GET /alarms
fn list(db: Db) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("alarms")
        .and(warp::get())
        .and(with_db(db))
        .and_then(list_alarms)
}

async fn list_alarms(db: Db) -> Result<impl warp::Reply, Infallible> {
    let alarms = get_alarms(db);
    Ok(warp::reply::json(&alarms))
}

fn get_alarms(db: Db) -> Vec<Alarm> {
    let mut alarms = vec![];
    db.read(|db| {
        alarms = db.alarms.values().cloned().collect();
    })
    .expect("Error while listing alarms");
    return alarms;
}

/// POST /alarms with JSON body
fn create(
    db: Db,
    mqtt_client: AsyncClient,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("alarms")
        .and(warp::post())
        .and(json_body())
        .and(with_db(db))
        .and(with_mqtt(mqtt_client))
        .and_then(create_alarm)
}

async fn create_alarm(
    alarm: Alarm,
    db: Db,
    mqtt_client: AsyncClient,
) -> Result<impl warp::Reply, Infallible> {
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
    mqtt::publish_alarms_changed(mqtt_client, get_alarms(db)).await;
    let json_reply = warp::reply::json(&alarm);
    Ok(warp::reply::with_status(json_reply, StatusCode::CREATED))
}

/// PUT /alarms/:id with JSON body
fn update(
    db: Db,
    mqtt_client: AsyncClient,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("alarms" / Uuid)
        .and(warp::put())
        .and(json_body())
        .and(with_db(db))
        .and(with_mqtt(mqtt_client))
        .and_then(update_alarm)
}

async fn update_alarm(
    id: Uuid,
    alarm: Alarm,
    db: Db,
    mqtt_client: AsyncClient,
) -> Result<impl warp::Reply, Infallible> {
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
    mqtt::publish_alarms_changed(mqtt_client, get_alarms(db)).await;
    let json_reply = warp::reply::json(&alarm);
    return Ok(warp::reply::with_status(json_reply, StatusCode::OK));
}

/// DELETE /alarms/:id
fn delete(
    db: Db,
    mqtt_client: AsyncClient,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("alarms" / Uuid)
        .and(warp::delete())
        .and(with_db(db))
        .and(with_mqtt(mqtt_client))
        .and_then(delete_alarm)
}

pub async fn delete_alarm(
    id: Uuid,
    db: Db,
    mqtt_client: AsyncClient,
) -> Result<impl warp::Reply, Infallible> {
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
        mqtt::publish_alarms_changed(mqtt_client, get_alarms(db)).await;
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

fn with_mqtt(
    mqtt_client: AsyncClient,
) -> impl Filter<Extract = (AsyncClient,), Error = std::convert::Infallible> + Clone {
    warp::any().map(move || mqtt_client.clone())
}

fn json_body() -> impl Filter<Extract = (Alarm,), Error = warp::Rejection> + Clone {
    // When accepting a body, we want a JSON body
    // (and to reject huge payloads)...
    warp::body::content_length_limit(1024 * 16).and(warp::body::json())
}
