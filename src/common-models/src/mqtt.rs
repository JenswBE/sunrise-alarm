use serde::{Deserialize, Serialize};

use crate::general::Alarm;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AlarmsChanged {
    pub alarms: Vec<Alarm>,
}
