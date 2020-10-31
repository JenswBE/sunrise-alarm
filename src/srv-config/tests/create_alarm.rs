use srv_config::{api, models};
use warp::http::StatusCode;
use warp::test::request;

mod common;
use common::fixture_alarm;

#[tokio::test]
async fn test_post() {
    let db = models::blank_db();
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
    let db = models::blank_db();
    let alarm = fixture_alarm();
    db.lock().await.push(alarm.clone());
    let api = api::alarms::filters(db);

    let resp = request()
        .method("POST")
        .path("/alarms")
        .json(&alarm)
        .reply(&api)
        .await;

    assert_eq!(resp.status(), StatusCode::CONFLICT);
}
