use srv_config::api;
use warp::http::StatusCode;
use warp::test::request;

mod common;
use common::*;

#[tokio::test]
async fn test_delete_alarm_success() {
    // Setup test data
    let db = fixture_database();
    let alarm = fixture_alarm();
    db.write(|db| db.alarms.insert(alarm.id, alarm.clone()))
        .unwrap();
    let api = api::alarms::filters(db);

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
}

#[tokio::test]
async fn test_delete_alarm_not_found() {
    // Setup test data
    let db = fixture_database();
    let api = api::alarms::filters(db);

    // Call service
    let resp = request()
        .method("DELETE")
        .path("/alarms/1")
        .json(&fixture_alarm())
        .reply(&api)
        .await;

    // Assert results
    assert_eq!(StatusCode::NOT_FOUND, resp.status());
}
