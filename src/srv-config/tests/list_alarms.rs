use warp::http::StatusCode;
use warp::test::request;

use srv_config::api;
use srv_config::Alarm;

mod common;
use common::*;

#[tokio::test]
async fn test_list_alarms_empty() {
    // Setup test data
    let db = fixture_database();
    let api = api::alarms::filters(db);

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
    let alarm = fixture_alarm();
    let alarm2 = fixture_alarm();
    db.write(|db| {
        db.alarms.insert(alarm.id, alarm.clone());
        db.alarms.insert(alarm2.id, alarm2.clone());
    })
    .unwrap();
    let api = api::alarms::filters(db);

    // Call service
    let resp = request().method("GET").path("/alarms").reply(&api).await;

    // Assert results
    assert_eq!(StatusCode::OK, resp.status());
    let result: Vec<Alarm> = serde_json::from_slice(resp.body()).unwrap();
    assert_eq!(2, result.len())
}
