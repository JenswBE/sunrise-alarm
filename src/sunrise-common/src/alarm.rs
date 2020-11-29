use chrono::{DateTime, Local};
use serde::{Deserialize, Serialize, Serializer};
use uuid::Uuid;

#[derive(Debug, Deserialize, Serialize, Default, Clone, Eq, PartialEq)]
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

#[derive(Debug, Serialize, Clone, Default)]
pub struct NextAlarms {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ring: Option<NextAlarm>,

    #[serde(skip_serializing_if = "Option::is_none")]
    pub action: Option<NextAlarm>,
}

#[derive(Debug, Serialize, Clone)]
pub struct NextAlarm {
    pub id: Uuid,
    pub alarm_datetime: DateTime<Local>,
    pub next_action: NextAction,
    pub next_action_datetime: DateTime<Local>,
}

#[derive(Debug, Clone, Eq, PartialEq)]
pub enum NextAction {
    None,
    Ring,
    Skip,
}

impl Default for NextAction {
    fn default() -> Self {
        NextAction::None
    }
}

impl Serialize for NextAction {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        match *self {
            NextAction::None => serializer.serialize_unit_variant("NextAction", 0, "None"),
            NextAction::Ring => serializer.serialize_unit_variant("NextAction", 1, "Ring"),
            NextAction::Skip => serializer.serialize_unit_variant("NextAction", 2, "Skip"),
        }
    }
}
