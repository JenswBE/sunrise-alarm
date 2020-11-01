#![deny(warnings)]

use std::env;

use srv_config::models;

#[tokio::main]
async fn main() {
    // Parse env variables
    let port = env::var("WARP_PORT")
        .unwrap_or("8000".to_string())
        .parse()
        .expect("Provided WARP_PORT is not a valid number");
    let data_dir = env::var("DATA_DIR_PATH")
        .unwrap_or("../../data".to_string())
        .parse()
        .expect("test");

    // Build config
    let config = models::Config { port, data_dir };

    // Run service
    srv_config::run(config).await;
}
