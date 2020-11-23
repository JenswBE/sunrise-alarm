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

#[derive(Debug, Clone, Eq, PartialEq)]
pub enum NextAction {
    None,
    Ring,
    Skip,
}

#[derive(Debug, Serialize, Clone, Default)]
pub struct NextAlarm {
    #[serde(skip_serializing_if = "Uuid::is_nil")]
    pub id: Uuid,

    #[serde(skip_serializing_if = "Option::is_some")]
    pub alarm_datetime: Option<DateTime<Local>>,

    #[serde(skip_serializing_if = "NextAction::is_none")]
    pub next_action: NextAction,

    #[serde(skip_serializing_if = "Option::is_some")]
    pub next_action_datetime: Option<DateTime<Local>>,
}

impl NextAction {
    fn is_none(&self) -> bool {
        *self == NextAction::None
    }
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
