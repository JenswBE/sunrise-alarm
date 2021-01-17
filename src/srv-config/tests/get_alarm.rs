use warp::http::StatusCode;
use warp::test::request;

use srv_config::api;
use sunrise_common::alarm::Alarm;

mod fixtures;
use fixtures::*;

#[tokio::test]
async fn test_get_alarm() {
    // Setup test data
    let db = fixture_database();
    let mqtt_client = fixture_mqtt_client().await;
    let alarm = fixture_alarm();
    db.write(|db| {
        db.alarms.insert(alarm.id, alarm.clone());
    })
    .unwrap();
    let api = api::alarms::filters(db, mqtt_client);

    // Call service
    let path = format!("/alarms/{}", alarm.id);
    let resp = request().method("GET").path(&path).reply(&api).await;

    // Assert results
    assert_eq!(StatusCode::OK, resp.status());
    let result: Alarm = serde_json::from_slice(resp.body()).unwrap();
    assert_eq!(alarm, result);
}

#[tokio::test]
async fn test_get_alarm_not_found() {
    // Setup test data
    let db = fixture_database();
    let mqtt_client = fixture_mqtt_client().await;
    let api = api::alarms::filters(db, mqtt_client);

    // Call service
    let path = "/alarms/00d60494-ad8c-44b7-af57-8c662271bd8a";
    let resp = request().method("GET").path(path).reply(&api).await;

    // Assert results
    assert_eq!(StatusCode::NOT_FOUND, resp.status());
}
