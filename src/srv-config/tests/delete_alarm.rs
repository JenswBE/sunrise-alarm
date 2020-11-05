use warp::http::StatusCode;
use warp::test::request;

use srv_config::api;

mod common;
use common::*;

#[tokio::test]
async fn test_delete_alarm_success() {
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
        .method("DELETE")
        .path(&path)
        .json(&alarm)
        .reply(&api)
        .await;

    // Assert results
    assert_eq!(StatusCode::NO_CONTENT, resp.status());

    // Assert MQTT server
    let stats = mqtt_server.fut.await.unwrap();
    assert_eq!(1, stats.len());
}

#[tokio::test]
async fn test_delete_alarm_not_found() {
    // Setup test data
    let db = fixture_database();
    let (mqtt_client, mqtt_server) = fixture_mqtt_client().await;
    let api = api::alarms::filters(db, mqtt_client);

    // Call service
    let resp = request()
        .method("DELETE")
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
