#![deny(warnings)]

use std::env;

use srv_config::models;
use sunrise_common::mqtt::MqttConfig;

#[tokio::main]
async fn main() {
    // Parse env variables
    let port = env::var("WARP_PORT")
        .unwrap_or("8001".to_string())
        .parse()
        .expect("Provided WARP_PORT is not a valid number");
    let data_dir = env::var("DATA_DIR_PATH")
        .unwrap_or("../../data".to_string())
        .parse()
        .expect("Provided DATA_DIR_PATH is not a valid path");
    let mqtt_broker_host = env::var("MQTT_BROKER_HOST").unwrap_or("localhost".to_string());
    let mqtt_broker_port = env::var("MQTT_BROKER_PORT")
        .unwrap_or("1883".to_string())
        .parse()
        .expect("Provided MQTT_BROKER_PORT is not a valid number");

    // Build config
    let config = models::Config {
        port,
        data_dir,
        mqtt_config: MqttConfig {
            host: mqtt_broker_host,
            port: mqtt_broker_port,
        },
    };

    // Run service
    srv_config::run(config).await;
}
