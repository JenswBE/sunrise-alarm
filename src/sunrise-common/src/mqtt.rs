use std::future::Future;

use rumqttc::{AsyncClient, MqttOptions, QoS};
use serde::{Deserialize, Serialize};

use crate::alarm::Alarm;

// =============================================
// =                 CONSTANTS                 =
// =============================================
const TOPIC_PREFIX: &str = "sunrise_alarm/";

// =============================================
// =                   TOPICS                  =
// =============================================
const TOPIC_ALARMS_CHANGED: &str = "alarms_changed";
const TOPIC_BUTTON_PRESSED: &str = "button_pressed";
const TOPIC_BUTTON_LONG_PRESSED: &str = "button_long_pressed";

// =============================================
// =                   CLIENT                  =
// =============================================
#[derive(Debug, Clone, Default)]
pub struct MqttConfig {
    pub host: String,
    pub port: u16,
}

pub async fn get_client<Fut>(
    config: MqttConfig,
    notification_handler: fn(rumqttc::Event) -> Fut,
) -> AsyncClient
where
    Fut: Future<Output = ()> + Send + 'static,
{
    // Build client
    let mut options = MqttOptions::new("srv-config", config.host, config.port);
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

// =============================================
// =           TOPIC: ALARMS CHANGED           =
// =============================================
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AlarmsChanged {
    pub alarms: Vec<Alarm>,
}

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
