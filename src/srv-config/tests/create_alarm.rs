use warp::http::StatusCode;
use warp::test::request;

use common_models::general::Alarm;
use srv_config::api;

mod common;
use common::*;

#[tokio::test]
async fn test_create_alarm_success() {
    // Setup test data
    let db = fixture_database();
    let (mqtt_client, mqtt_server) = fixture_mqtt_client().await;
    let api = api::alarms::filters(db, mqtt_client);
    let alarm = &fixture_alarm();

    // Call service
    let resp = request()
        .method("POST")
        .path("/alarms")
        .json(&alarm)
        .reply(&api)
        .await;

    // Assert results
    assert_eq!(StatusCode::CREATED, resp.status());
    let result: Alarm = serde_json::from_slice(resp.body()).unwrap();
    assert_eq!(*alarm, result);

    // Assert MQTT server
    let stats = mqtt_server.fut.await.unwrap();
    assert_eq!(1, stats.len());
}

#[tokio::test]
async fn test_create_alarm_conflict() {
    // Setup test data
    let db = fixture_database();
    let (mqtt_client, mqtt_server) = fixture_mqtt_client().await;
    let alarm = fixture_alarm();
    db.write(|db| db.alarms.insert(alarm.id, alarm.clone()))
        .unwrap();
    let api = api::alarms::filters(db, mqtt_client);

    // Call service
    let resp = request()
        .method("POST")
        .path("/alarms")
        .json(&alarm)
        .reply(&api)
        .await;

    // Assert results
    assert_eq!(StatusCode::CONFLICT, resp.status());

    // Assert MQTT server
    let stats = mqtt_server.fut.await.unwrap();
    assert_eq!(1, stats.len());
}
