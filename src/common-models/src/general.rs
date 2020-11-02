use serde::{Deserialize, Serialize};
use uuid::Uuid;

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
