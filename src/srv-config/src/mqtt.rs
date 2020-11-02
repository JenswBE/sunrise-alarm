use mqtt_async_client::client::{Client, Publish, QoS};

use crate::models;

pub async fn publish(config: models::MqttConfig, topic: &str, message: Vec<u8>) {
    let mut client = get_client(config);
    client
        .connect()
        .await
        .expect("Failed to connect to MQTT broker");
    let mut p = Publish::new("sunrise_alarm/".to_string() + topic, message);
    p.set_qos(QoS::AtLeastOnce);
    client.publish(&p).await.expect("Failed to publish message")
}

fn get_client(config: models::MqttConfig) -> Client {
    let mut b = Client::builder();
    b.set_host(config.host)
        .set_port(config.port)
        .set_client_id(Some("srv-config".to_string()))
        .build()
        .expect("Failed to create MQTT client")
}
