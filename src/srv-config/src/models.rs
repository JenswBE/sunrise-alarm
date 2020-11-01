use std::collections::HashMap;
use std::path::PathBuf;

use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Clone)]
pub struct Config {
    pub data_dir: PathBuf,
    pub port: u16,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ServerData {
    pub alarms: HashMap<Uuid, Alarm>,
}

#[derive(Debug, Deserialize, Serialize, Clone, Eq, PartialEq)]
pub struct Alarm {
    #[serde(default = "Uuid::new_v4")]
    pub id: Uuid,
    pub enabled: bool,
    pub name: String,
    pub hour: u8,
    pub minute: u8,
    pub days: Vec<u8>,
    pub skip_next: bool,
}

pub const ERROR_ALARM_EXISTS: &str = "alarm_exists";
pub const ERROR_ALARM_NOT_FOUND: &str = "alarm_not_found";

#[derive(Debug, Serialize, Clone)]
pub struct Error {
    pub code: &'static str,
}
