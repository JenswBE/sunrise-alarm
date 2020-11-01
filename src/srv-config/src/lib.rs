#![deny(warnings)]

use std::env;

use warp::Filter;

pub mod api;
pub mod database;
pub mod models;

pub use models::Alarm;

/// Provides a RESTful web server for managing Sunrise Alarm's config
///
/// API contains following routes:
///
/// - `GET /alarms`: return a JSON list of Alarms.
/// - `POST /alarms`: create a new Alarm.
/// - `PUT /alarms/:id`: update a specific Alarm.
/// - `DELETE /alarms/:id`: delete a specific Alarm.
pub async fn run(config: models::Config) {
    // Setup logging
    if env::var_os("RUST_LOG").is_none() {
        // Set `RUST_LOG=alarms=debug` to see debug logs
        env::set_var("RUST_LOG", "alarms=info");
    }
    pretty_env_logger::init();

    // Setup database
    let mut db_path = config.data_dir.clone();
    db_path.push("server_data.ron");
    let db = database::load_or_init(db_path);
    let api = api::alarms::filters(db);

    // Add middleware for access logs
    let routes = api.with(warp::log("alarms"));

    // Start up the server
    warp::serve(routes).run(([127, 0, 0, 1], config.port)).await;
}
