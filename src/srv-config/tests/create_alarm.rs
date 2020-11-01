use srv_config::api;
use warp::http::StatusCode;
use warp::test::request;

mod common;
use common::*;

#[tokio::test]
async fn test_post() {
    let db = fixture_database();
    let api = api::alarms::filters(db);

    let resp = request()
        .method("POST")
        .path("/alarms")
        .json(&fixture_alarm())
        .reply(&api)
        .await;

    assert_eq!(resp.status(), StatusCode::CREATED);
}

#[tokio::test]
async fn test_post_conflict() {
    let db = fixture_database();
    let alarm = fixture_alarm();
    db.write(|db| db.alarms.insert(alarm.id, alarm.clone()))
        .unwrap();
    let api = api::alarms::filters(db);

    let resp = request()
        .method("POST")
        .path("/alarms")
        .json(&alarm)
        .reply(&api)
        .await;

    assert_eq!(resp.status(), StatusCode::CONFLICT);
}
