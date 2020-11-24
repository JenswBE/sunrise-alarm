#![deny(warnings)]

use srv_config::models::Config;

#[tokio::main]
async fn main() {
    // Build config
    let config = Config::from_env();

    // Run service
    srv_config::run(config).await;
}
