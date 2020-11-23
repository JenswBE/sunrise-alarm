use rumqttc::AsyncClient;

use sunrise_common::mqtt::MqttConfig;

pub async fn get_client(config: MqttConfig) -> AsyncClient {
    sunrise_common::mqtt::get_client(config, handle_mqtt_notification).await
}

async fn handle_mqtt_notification(_notification: rumqttc::Event) {}
