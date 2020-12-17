use std::collections::HashMap;
use std::sync::{Arc, Mutex};

use chrono::{DateTime, Duration, Local};
use reqwest::Client;
use rumqttc::AsyncClient;
use uuid::Uuid;

use crate::manager::Radio;
use sunrise_common::alarm::{Alarm, NextAlarm, NextAlarms};
use sunrise_common::config::{HostsConfig, MqttConfig, WarpConfig};

#[derive(Debug, Clone)]
pub struct Context {
    pub client: Client,
    pub mqtt_client: Arc<Mutex<Option<AsyncClient>>>,
    pub config: Arc<Config>,
    pub radio: Radio,
    state: Arc<Mutex<State>>,
}

impl Context {
    pub fn new(config: Config, radio: Radio) -> Self {
        Self {
            client: reqwest::Client::new(),
            mqtt_client: Arc::new(Mutex::new(None)),
            config: Arc::new(config),
            radio,
            state: Arc::new(Mutex::new(State::default())),
        }
    }

    pub fn get_mqtt_client(&self) -> AsyncClient {
        let client = self.mqtt_client.lock().unwrap();
        client.clone().unwrap()
    }

    pub fn set_mqtt_client(&self, new_client: AsyncClient) {
        let mut client = self.mqtt_client.lock().unwrap();
        *client = Some(new_client);
    }

    pub fn get_status(&self) -> Status {
        let state = self.state.lock().unwrap();
        state.status.clone()
    }

    pub fn set_status(&self, status: Status) {
        let mut state = self.state.lock().unwrap();
        state.status = status;
    }

    pub fn get_alarm(&self, id: Uuid) -> Option<Alarm> {
        let state = self.state.lock().unwrap();
        state.alarms.iter().find(|&a| a.id == id).cloned()
    }

    pub fn get_alarms(&self) -> Vec<Alarm> {
        let state = self.state.lock().unwrap();
        state.alarms.clone()
    }

    pub fn set_alarms(&self, alarms: Vec<Alarm>) {
        let mut state = self.state.lock().unwrap();
        state.alarms = alarms;
    }

    pub fn get_next_alarms(&self) -> NextAlarms {
        let state = self.state.lock().unwrap();
        state.next_alarms.clone()
    }

    pub fn get_next_alarms_ring(&self) -> Option<NextAlarm> {
        let state = self.state.lock().unwrap();
        state.next_alarms.ring.clone()
    }

    pub fn set_next_alarms_ring(&self, next_alarm: Option<NextAlarm>) {
        let mut state = self.state.lock().unwrap();
        state.next_alarms.ring = next_alarm;
    }

    pub fn get_next_alarms_action(&self) -> Option<NextAlarm> {
        let state = self.state.lock().unwrap();
        state.next_alarms.action.clone()
    }

    pub fn set_next_alarms_action(&self, next_alarm: Option<NextAlarm>) {
        let mut state = self.state.lock().unwrap();
        state.next_alarms.action = next_alarm;
    }

    pub fn get_last_ring(&self, alarm_id: Uuid) -> Option<DateTime<Local>> {
        let state = self.state.lock().unwrap();
        state.last_rings.get(&alarm_id).cloned()
    }

    pub fn set_last_ring(&self, alarm_id: Uuid, last_ring: DateTime<Local>) {
        let mut state = self.state.lock().unwrap();
        state.last_rings.insert(alarm_id, last_ring);
    }
}

#[derive(Debug, Clone)]
pub struct Config {
    pub alarm: AlarmConfig,
    pub hosts: HostsConfig,
    pub mqtt: MqttConfig,
    pub warp: WarpConfig,
}

impl Config {
    pub fn from_env() -> Self {
        Self {
            alarm: AlarmConfig::default(),
            hosts: HostsConfig::from_env(),
            mqtt: MqttConfig::from_env(),
            warp: WarpConfig::from_env(8000),
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
        Self {
            light_duration: Duration::minutes(10),
            sound_duration: Duration::minutes(5),
        }
    }
}

#[derive(Debug, Default)]
pub struct State {
    pub status: Status,
    pub alarms: Vec<Alarm>,
    pub next_alarms: NextAlarms,
    pub last_rings: HashMap<Uuid, DateTime<Local>>,
}

#[derive(Debug, Clone, Eq, PartialEq)]
pub enum Status {
    Idle,
    Ring(Uuid),
}

impl Default for Status {
    fn default() -> Self {
        Self::Idle
    }
}
