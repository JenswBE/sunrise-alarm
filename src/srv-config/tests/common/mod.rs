use std::collections::HashMap;
use std::env;
use std::sync::Arc;

use rustbreak::PathDatabase;
use tempfile::NamedTempFile;
use uuid::Uuid;

use common_models::general::Alarm;
use srv_config::database;
use srv_config::models;

pub fn fixture_alarm() -> Alarm {
    Alarm {
        id: Uuid::new_v4(),
        enabled: true,
        name: "New name".to_string(),
        hour: 12,
        minute: 0,
        days: vec![],
        skip_next: false,
    }
}

pub fn fixture_database() -> database::Db {
    let file = NamedTempFile::new().unwrap();
    let path = file.into_temp_path().to_path_buf();
    let db = PathDatabase::load_from_path_or(
        path,
        models::ServerData {
            alarms: HashMap::new(),
        },
    )
    .unwrap();
    Arc::new(db)
}

pub fn fixture_mqtt_config() -> models::MqttConfig {
    models::MqttConfig {
        host: env::var("MQTT_BROKER_HOST").unwrap_or("localhost".to_string()),
        port: 1883,
    }
}
