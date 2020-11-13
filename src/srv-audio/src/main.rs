#![deny(warnings)]

use std::env;

use srv_audio::models;
use sunrise_common::mqtt::MqttConfig;

#[tokio::main]
async fn main() {
    // Parse env variables
    let warp_port = env::var("WARP_PORT")
        .unwrap_or("8003".to_string())
        .parse()
        .expect("Provided WARP_PORT is not a valid number");
    let music_dir = env::var("MUSIC_DIR_PATH")
        .unwrap_or("../../data/music".to_string())
        .parse()
        .expect("Provided MUSIC_DIR_PATH is not a valid path");
    let mqtt_broker_host = env::var("MQTT_BROKER_HOST").unwrap_or("localhost".to_string());
    let mqtt_broker_port = env::var("MQTT_BROKER_PORT")
        .unwrap_or("1883".to_string())
        .parse()
        .expect("Provided MQTT_BROKER_PORT is not a valid number");

    // Build config
    let config = models::Config {
        warp_port,
        music_dir,
        mqtt_config: MqttConfig {
            host: mqtt_broker_host,
            port: mqtt_broker_port,
        },
    };

    // Run service
    srv_audio::run(config).await;
}
