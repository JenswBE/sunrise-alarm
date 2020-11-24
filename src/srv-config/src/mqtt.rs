use rumqttc::{AsyncClient, MqttOptions};

use sunrise_common::config::MqttConfig;

const CLIENT_ID: &str = "srv-config";

pub async fn get_client(config: MqttConfig) -> AsyncClient {
    // Build client
    let mut options = MqttOptions::new(CLIENT_ID, config.host, config.port);
    options.set_keep_alive(5);
    let (mqtt_client, mut eventloop) = AsyncClient::new(options, 10);

    // Start client loop
    tokio::task::spawn(async move {
        loop {
            let notification = eventloop.poll().await.unwrap();
            notification_handler(notification).await;
        }
    });

    // Create client successful
    return mqtt_client;
}

async fn notification_handler(_notification: rumqttc::Event) {}
