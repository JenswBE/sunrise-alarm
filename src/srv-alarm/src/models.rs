use serde::Serialize;
use std::sync::{Arc, Mutex};
use uuid::Uuid;

use sunrise_common::mqtt::MqttConfig;

#[derive(Debug, Clone)]
pub struct Config {
    pub port: u16,
    pub mqtt_config: MqttConfig,
}

pub type State = Arc<Mutex<LocalState>>;

impl LocalState {
    pub fn new() -> State {
        Arc::new(Mutex::new(LocalState::default()))
    }
}

#[derive(Debug, Clone, Default)]
pub struct LocalState {
    pub next_alarm: Option<NextAlarm>,
}

#[derive(Debug, Serialize, Clone, Default)]
pub struct NextAlarm {
    #[serde(skip_serializing_if = "Uuid::is_nil")]
    pub id: Uuid,

    #[serde(skip_serializing_if = "String::is_empty")]
    pub alarm_datetime: String,

    #[serde(skip_serializing_if = "String::is_empty")]
    pub next_action: String,

    #[serde(skip_serializing_if = "String::is_empty")]
    pub next_action_datetime: String,
}

#[derive(Debug, Serialize, Clone)]
pub struct Error {
    pub code: &'static str,
}
