#![deny(warnings)]

use std::env;
use warp::Filter;

pub use models::Alarm;

pub mod api;
pub mod models;

/// Provides a RESTful web server for managing Sunrise Alarm's config
///
/// API will be:
///
/// - `GET /alarms`: return a JSON list of Alarms.
/// - `POST /alarms`: create a new Alarm.
/// - `PUT /alarms/:id`: update a specific Alarm.
/// - `DELETE /alarms/:id`: delete a specific Alarm.
pub async fn run() {
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
    let port = env::var("WARP_PORT").unwrap_or("8000".to_string());
    let port = port
        .parse()
        .expect("Provided WARP_PORT is not a valid number");
    warp::serve(routes).run(([127, 0, 0, 1], port)).await;
}
