use std::collections::HashMap;
use std::path::PathBuf;

use serde::{Deserialize, Serialize};
use uuid::Uuid;

use sunrise_common::alarm::Alarm;
use sunrise_common::config::{parse_path, MqttConfig, WarpConfig};

#[derive(Debug, Clone)]
pub struct Config {
    pub data_dir: PathBuf,
    pub warp: WarpConfig,
    pub mqtt: MqttConfig,
}

impl Config {
    pub fn from_env() -> Self {
        Config {
            data_dir: parse_path("DATA_DIR_PATH", "../../data"),
            warp: WarpConfig::from_env(8001),
            mqtt: MqttConfig::from_env(),
        }
    }
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
