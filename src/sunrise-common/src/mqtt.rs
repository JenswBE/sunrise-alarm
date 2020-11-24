use rumqttc::{AsyncClient, Publish, QoS};
use serde::{Deserialize, Serialize};

use crate::alarm::Alarm;

// =============================================
// =                   TOPICS                  =
// =============================================
pub const TOPIC_ALARMS_CHANGED: &str = "sunrise_alarm/alarms_changed";
pub const TOPIC_BUTTON_PRESSED: &str = "sunrise_alarm/button_pressed";
pub const TOPIC_BUTTON_LONG_PRESSED: &str = "sunrise_alarm/button_long_pressed";

// =============================================
// =                   CLIENT                  =
// =============================================
pub async fn subscribe(client: &AsyncClient, topic: &str) {
    client.subscribe(topic, QoS::AtLeastOnce).await.unwrap()
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
        .publish(TOPIC_ALARMS_CHANGED, QoS::AtLeastOnce, false, json)
        .await
        .unwrap();
}

pub fn parse_alarms_changed(packet: Publish) -> AlarmsChanged {
    serde_json::from_slice(&packet.payload[..]).unwrap()
}
