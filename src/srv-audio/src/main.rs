#![deny(warnings)]

use srv_audio::models::Config;

#[tokio::main]
async fn main() {
    // Build config
    let config = Config::from_env();

    // Run service
    srv_audio::run(config).await;
}
