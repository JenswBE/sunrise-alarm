#![deny(warnings)]
#![deny(missing_debug_implementations)]

use std::env;

use warp::Filter;

pub mod api;
pub mod http;
pub mod manager;
pub mod models;
pub mod mqtt;
pub mod planner;
pub mod ringer;

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

    // Create context
    let (radio, receiver) = manager::create_radios();
    let ctx = models::Context::new(config, radio);

    // Fetch alarms
    let alarms = http::get_alarms(&ctx)
        .await
        .map_err(|_| std::process::exit(1))
        .unwrap();
    ctx.set_alarms(alarms);

    // Setup manager
    manager::start(ctx.clone(), receiver);

    // Setup MQTT
    let _mqtt_client = mqtt::get_client(ctx.clone()).await;

    // Initial update of next alarms
    planner::update_next_alarms(&ctx).await;

    // Setup server
    let api = api::alarms::filters(ctx.clone());
    let routes = api.with(warp::log("alarm"));

    // Start the server
    warp::serve(routes)
        .run(([0, 0, 0, 0], ctx.config.warp.port))
        .await;
}
