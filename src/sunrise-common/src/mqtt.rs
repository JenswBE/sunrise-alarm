use serde::{Deserialize, Serialize};

use crate::general::Alarm;

// Topics
const TOPIC_ALARMS_CHANGED: &str = "alarms_changed";
const TOPIC_BUTTON_PRESSED: &str = "button_pressed";
const TOPIC_BUTTON_LONG_PRESSED: &str = "button_long_pressed";

// Models
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AlarmsChanged {
    pub alarms: Vec<Alarm>,
}
