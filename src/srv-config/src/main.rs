#![deny(warnings)]

use std::env;
use warp::Filter;

mod api;
mod models;

/// Provides a RESTful web server for managing Sunrise Alarm's config
///
/// API will be:
///
/// - `GET /alarms`: return a JSON list of Alarms.
/// - `POST /alarms`: create a new Alarm.
/// - `PUT /alarms/:id`: update a specific Alarm.
/// - `DELETE /alarms/:id`: delete a specific Alarm.
#[tokio::main]
async fn main() {
    if env::var_os("RUST_LOG").is_none() {
        // Set `RUST_LOG=alarms=debug` to see debug logs,
        // this only shows access logs.
        env::set_var("RUST_LOG", "alarms=info");
    }
    pretty_env_logger::init();

    let db = models::blank_db();

    let api = api::alarms::filters(db);

    // View access logs by setting `RUST_LOG=alarms`.
    let routes = api.with(warp::log("alarms"));
    // Start up the server...
    warp::serve(routes).run(([127, 0, 0, 1], 3030)).await;
}

// #[cfg(test)]
// mod tests {
//     use warp::http::StatusCode;
//     use warp::test::request;

//     use super::{
//         filters,
//         models::{self, Alarm},
//     };

//     #[tokio::test]
//     async fn test_post() {
//         let db = models::blank_db();
//         let api = filters::alarms(db);

//         let resp = request()
//             .method("POST")
//             .path("/alarms")
//             .json(&Alarm {
//                 id: 1,
//                 text: "test 1".into(),
//                 completed: false,
//             })
//             .reply(&api)
//             .await;

//         assert_eq!(resp.status(), StatusCode::CREATED);
//     }

//     #[tokio::test]
//     async fn test_post_conflict() {
//         let db = models::blank_db();
//         db.lock().await.push(alarm1());
//         let api = filters::alarms(db);

//         let resp = request()
//             .method("POST")
//             .path("/alarms")
//             .json(&alarm1())
//             .reply(&api)
//             .await;

//         assert_eq!(resp.status(), StatusCode::BAD_REQUEST);
//     }

//     #[tokio::test]
//     async fn test_put_unknown() {
//         let _ = pretty_env_logger::try_init();
//         let db = models::blank_db();
//         let api = filters::alarms(db);

//         let resp = request()
//             .method("PUT")
//             .path("/alarms/1")
//             .header("authorization", "Bearer admin")
//             .json(&alarm1())
//             .reply(&api)
//             .await;

//         assert_eq!(resp.status(), StatusCode::NOT_FOUND);
//     }

//     fn alarm1() -> Alarm {
//         Alarm {
//             id: 1,
//             text: "test 1".into(),
//             completed: false,
//         }
//     }
// }
