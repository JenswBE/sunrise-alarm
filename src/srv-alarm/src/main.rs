#![deny(warnings)]

use chrono::Duration;
use std::env;

use srv_alarm::models;
use sunrise_common::mqtt::MqttConfig;

#[tokio::main]
async fn main() {
    // Parse env variables
    let port = env::var("WARP_PORT")
        .unwrap_or("8001".to_string())
        .parse()
        .expect("Provided WARP_PORT is not a valid number");
    let mqtt_broker_host = env::var("MQTT_BROKER_HOST").unwrap_or("localhost".to_string());
    let mqtt_broker_port = env::var("MQTT_BROKER_PORT")
        .unwrap_or("1883".to_string())
        .parse()
        .expect("Provided MQTT_BROKER_PORT is not a valid number");

    // Build config
    let config = models::Config {
        port,
        alarm_config: models::AlarmConfig {
            light_duration: Duration::minutes(10),
            sound_duration: Duration::minutes(5),
        },
        mqtt_config: MqttConfig {
            host: mqtt_broker_host,
            port: mqtt_broker_port,
        },
    };

    // Run service
    srv_alarm::run(config).await;
}
