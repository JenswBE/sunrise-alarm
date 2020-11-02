use warp::http::StatusCode;
use warp::test::request;

use common_models::general::Alarm;
use srv_config::api;

mod common;
use common::*;

#[tokio::test]
async fn test_update_alarm_success() {
    // Setup test data
    let db = fixture_database();
    let mqtt_config = fixture_mqtt_config();
    let alarm = fixture_alarm();
    db.write(|db| db.alarms.insert(alarm.id, alarm.clone()))
        .unwrap();
    let api = api::alarms::filters(db, mqtt_config);

    // Call service
    let path = format!("/alarms/{}", alarm.id);
    let resp = request()
        .method("PUT")
        .path(&path)
        .json(&alarm)
        .reply(&api)
        .await;

    // Assert results
    assert_eq!(StatusCode::OK, resp.status());
    let result: Alarm = serde_json::from_slice(resp.body()).unwrap();
    assert_eq!(alarm, result)
}

#[tokio::test]
async fn test_update_alarm_not_found() {
    // Setup test data
    let db = fixture_database();
    let mqtt_config = fixture_mqtt_config();
    let api = api::alarms::filters(db, mqtt_config);

    // Call service
    let resp = request()
        .method("PUT")
        .path("/alarms/1")
        .json(&fixture_alarm())
        .reply(&api)
        .await;

    // Assert results
    assert_eq!(StatusCode::NOT_FOUND, resp.status());
}
