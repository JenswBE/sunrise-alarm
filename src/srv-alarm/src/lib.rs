#![deny(warnings)]

use std::env;

use warp::Filter;

pub mod api;
pub mod models;
pub mod mqtt;

/// Provides a RESTful web server for general management of Sunrise Alarm
///
/// API contains following routes:
///
/// - `GET /alarms/next`: returns information about the next alarm.
pub async fn run(config: models::Config) {
    // Setup logging
    if env::var_os("RUST_LOG").is_none() {
        // Set `RUST_LOG=debug` to see debug logs
        env::set_var("RUST_LOG", "info");
    }
    pretty_env_logger::init();

    // Setup state
    let state = models::LocalState::new();

    // Setup MQTT
    let _mqtt_client = mqtt::get_client(config.mqtt_config).await;

    // Setup server
    let api = api::alarms::filters(state);
    let routes = api.with(warp::log("alarm"));

    // Start the server
    warp::serve(routes).run(([0, 0, 0, 0], config.port)).await;
}
