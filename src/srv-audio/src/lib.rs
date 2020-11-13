#![deny(warnings)]

use std::env;

use warp::Filter;

pub mod api;
pub mod models;
pub mod player;

pub async fn run(config: models::Config) {
    // Setup logging
    if env::var_os("RUST_LOG").is_none() {
        // Set `RUST_LOG=audio=debug` to see debug logs
        env::set_var("RUST_LOG", "audio=info");
    }
    pretty_env_logger::init();

    // Setup player
    let remote = player::setup().await;

    // Setup server
    let api = api::music::filters(remote);
    let routes = api.with(warp::log("audio"));

    // Start the server
    warp::serve(routes)
        .run(([0, 0, 0, 0], config.warp_port))
        .await;
}
