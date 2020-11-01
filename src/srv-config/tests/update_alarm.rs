use srv_config::{api, models};
use warp::http::StatusCode;
use warp::test::request;

mod common;
use common::fixture_alarm;

// #[tokio::test]
// async fn test_put_unknown() {
//     let _ = pretty_env_logger::try_init();
//     let db = models::blank_db();
//     let api = api::alarms::filters(db);

//     let resp = request()
//         .method("PUT")
//         .path("/alarms/1")
//         .json(&fixture_alarm())
//         .reply(&api)
//         .await;

//     assert_eq!(resp.status(), StatusCode::NOT_FOUND);
// }
