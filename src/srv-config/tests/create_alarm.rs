use warp::http::StatusCode;
use warp::test::request;

use srv_config::api;
use srv_config::Alarm;

mod common;
use common::*;

#[tokio::test]
async fn test_create_alarm_success() {
    // Setup test data
    let db = fixture_database();
    let api = api::alarms::filters(db);
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
    assert_eq!(*alarm, result)
}

#[tokio::test]
async fn test_create_alarm_conflict() {
    // Setup test data
    let db = fixture_database();
    let alarm = fixture_alarm();
    db.write(|db| db.alarms.insert(alarm.id, alarm.clone()))
        .unwrap();
    let api = api::alarms::filters(db);

    // Call service
    let resp = request()
        .method("POST")
        .path("/alarms")
        .json(&alarm)
        .reply(&api)
        .await;

    // Assert results
    assert_eq!(StatusCode::CONFLICT, resp.status());
}
