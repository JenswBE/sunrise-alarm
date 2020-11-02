use warp::http::StatusCode;
use warp::test::request;

use common_models::general::Alarm;
use srv_config::api;

mod common;
use common::*;

#[tokio::test]
async fn test_list_alarms_empty() {
    // Setup test data
    let db = fixture_database();
    let mqtt_config = fixture_mqtt_config();
    let api = api::alarms::filters(db, mqtt_config);

    // Call service
    let resp = request().method("GET").path("/alarms").reply(&api).await;

    // Assert results
    assert_eq!(StatusCode::OK, resp.status());
    let result: Vec<Alarm> = serde_json::from_slice(resp.body()).unwrap();
    assert_eq!(0, result.len())
}

#[tokio::test]
async fn test_list_alarms_not_empty() {
    // Setup test data
    let db = fixture_database();
    let mqtt_config = fixture_mqtt_config();
    let alarm = fixture_alarm();
    let alarm2 = fixture_alarm();
    db.write(|db| {
        db.alarms.insert(alarm.id, alarm.clone());
        db.alarms.insert(alarm2.id, alarm2.clone());
    })
    .unwrap();
    let api = api::alarms::filters(db, mqtt_config);

    // Call service
    let resp = request().method("GET").path("/alarms").reply(&api).await;

    // Assert results
    assert_eq!(StatusCode::OK, resp.status());
    let result: Vec<Alarm> = serde_json::from_slice(resp.body()).unwrap();
    assert_eq!(2, result.len())
}
