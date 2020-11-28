#![deny(warnings)]
#![deny(missing_debug_implementations)]

use srv_alarm::models::Config;

#[tokio::main]
async fn main() {
    // Build config
    let config = Config::from_env();

    // Run service
    srv_alarm::run(config).await;
}
