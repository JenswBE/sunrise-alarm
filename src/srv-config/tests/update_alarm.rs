use warp::http::StatusCode;
use warp::test::request;

use srv_config::api;
use sunrise_common::alarm::Alarm;

mod fixtures;
use fixtures::*;

#[tokio::test]
async fn test_update_alarm_success() {
    // Setup test data
    let db = fixture_database();
    let (mqtt_client, mqtt_server) = fixture_mqtt_client().await;
    let alarm = fixture_alarm();
    db.write(|db| db.alarms.insert(alarm.id, alarm.clone()))
        .unwrap();
    let api = api::alarms::filters(db, mqtt_client);

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
    assert_eq!(alarm, result);

    // Assert MQTT server
    let stats = mqtt_server.fut.await.unwrap();
    assert_eq!(1, stats.len());
}

#[tokio::test]
async fn test_update_alarm_not_found() {
    // Setup test data
    let db = fixture_database();
    let (mqtt_client, mqtt_server) = fixture_mqtt_client().await;
    let api = api::alarms::filters(db, mqtt_client);

    // Call service
    let resp = request()
        .method("PUT")
        .path("/alarms/1")
        .json(&fixture_alarm())
        .reply(&api)
        .await;

    // Assert results
    assert_eq!(StatusCode::NOT_FOUND, resp.status());

    // Assert MQTT server
    let stats = mqtt_server.fut.await.unwrap();
    assert_eq!(1, stats.len());
}
