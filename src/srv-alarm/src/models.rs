use chrono::Duration;
use std::sync::{Arc, Mutex};

use sunrise_common::alarm::{Alarm, NextAlarm};
use sunrise_common::mqtt::MqttConfig;

#[derive(Debug, Clone)]
pub struct Config {
    pub port: u16,
    pub light_duration: Duration,
    pub sound_duration: Duration,
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
    pub alarms: Vec<Alarm>,
    pub next_alarm: Option<NextAlarm>,
}
