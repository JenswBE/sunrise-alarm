use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::sync::Mutex;
use uuid::Uuid;

/// So we don't have to tackle how different database work, we'll just use
/// a simple in-memory DB, a vector synchronized by a mutex.
pub type Db = Arc<Mutex<Vec<Alarm>>>;

pub fn blank_db() -> Db {
    Arc::new(Mutex::new(Vec::new()))
}

#[derive(Debug, Serialize, Deserialize, Clone)]
struct ServerData {
    alarms: Vec<Alarm>,
}

#[derive(Debug, Deserialize, Serialize, Clone)]
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
