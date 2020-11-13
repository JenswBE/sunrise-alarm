use std::collections::HashMap;
use std::path::PathBuf;

use serde::{Deserialize, Serialize};
use uuid::Uuid;

use sunrise_common::general::Alarm;
use sunrise_common::mqtt::MqttConfig;

#[derive(Debug, Clone)]
pub struct Config {
    pub data_dir: PathBuf,
    pub port: u16,
    pub mqtt_config: MqttConfig,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ServerData {
    pub alarms: HashMap<Uuid, Alarm>,
}

pub const ERROR_ALARM_EXISTS: &str = "alarm_exists";
pub const ERROR_ALARM_NOT_FOUND: &str = "alarm_not_found";

#[derive(Debug, Serialize, Clone)]
pub struct Error {
    pub code: &'static str,
}
