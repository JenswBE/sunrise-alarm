use std::sync::{Arc, Mutex};

use chrono::Duration;

use sunrise_common::alarm::{Alarm, NextAlarm};
use sunrise_common::config::{MqttConfig, WarpConfig};

#[derive(Debug, Clone)]
pub struct Config {
    pub warp: WarpConfig,
    pub alarm: AlarmConfig,
    pub mqtt: MqttConfig,
}

impl Config {
    pub fn from_env() -> Self {
        Config {
            warp: WarpConfig::from_env(8000),
            alarm: AlarmConfig::default(),
            mqtt: MqttConfig::from_env(),
        }
    }
}

#[derive(Debug, Clone)]
pub struct AlarmConfig {
    pub light_duration: Duration,
    pub sound_duration: Duration,
}

impl Default for AlarmConfig {
    fn default() -> Self {
        AlarmConfig {
            light_duration: Duration::minutes(10),
            sound_duration: Duration::minutes(5),
        }
    }
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
    pub next_alarms: Vec<NextAlarm>,
    pub next_alarm_ring: Option<NextAlarm>,
    pub next_alarm_action: Option<NextAlarm>,
}
