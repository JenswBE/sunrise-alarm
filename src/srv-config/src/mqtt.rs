use rumqttc::{AsyncClient, MqttOptions, QoS};

use crate::models;
use sunrise_common::{general::Alarm, mqtt::AlarmsChanged};

const TOPIC_PREFIX: &str = "sunrise_alarm/";

pub async fn get_client(config: models::MqttConfig) -> AsyncClient {
    // Build client
    let mut mqttoptions = MqttOptions::new("srv-config", config.host, config.port);
    mqttoptions.set_keep_alive(5);
    let (mqtt_client, mut eventloop) = AsyncClient::new(mqttoptions, 10);

    // Start client loop
    tokio::task::spawn(async move {
        loop {
            let notification = eventloop.poll().await.unwrap();
            handle_mqtt_notification(notification).await;
        }
    });

    // Create client successful
    return mqtt_client;
}

async fn handle_mqtt_notification(_notification: rumqttc::Event) {}

pub async fn publish_alarms_changed(client: AsyncClient, alarms: Vec<Alarm>) {
    let msg = AlarmsChanged { alarms };
    let json = serde_json::to_vec(&msg).unwrap();
    client
        .publish(
            TOPIC_PREFIX.to_string() + "alarms_changed",
            QoS::AtLeastOnce,
            false,
            json,
        )
        .await
        .unwrap();
}
