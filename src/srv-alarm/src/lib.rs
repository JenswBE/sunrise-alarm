#![deny(warnings)]

use std::env;

use warp::Filter;

pub mod api;
pub mod manager;
pub mod models;
pub mod mqtt;
pub mod time;

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

    // Validate config
    if config.alarm.light_duration < config.alarm.sound_duration {
        panic!("Sound should start after or together with light (duration light >= sound)")
    }

    // Setup state
    let state = models::LocalState::new();
    {
        let mut state = state.lock().unwrap();
        state.alarms = reqwest::get("http://localhost:8001/alarms")
            .await
            .unwrap()
            .json()
            .await
            .unwrap();
    }

    // Setup manager
    let radio = manager::start(state.clone(), config.clone());

    // Initial update of next alarms
    time::update_next_alarms(state.clone(), &config.alarm, radio.clone());

    // Setup MQTT
    let _mqtt_client = mqtt::get_client(&config, state.clone(), radio).await;

    // Setup server
    let api = api::alarms::filters(state);
    let routes = api.with(warp::log("alarm"));

    // Start the server
    warp::serve(routes)
        .run(([0, 0, 0, 0], config.warp.port))
        .await;
}
