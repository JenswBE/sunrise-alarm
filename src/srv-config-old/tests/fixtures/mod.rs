use std::collections::HashMap;
use std::sync::Arc;

use rustbreak::PathDatabase;
use tempfile::NamedTempFile;
use uuid::Uuid;

use srv_config::database;
use srv_config::models;
use srv_config::mqtt;
use sunrise_common::alarm::Alarm;
use sunrise_common::config::MqttConfig;

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

pub async fn fixture_mqtt_client() -> rumqttc::AsyncClient {
    let config = MqttConfig {
        host: "localhost".to_string(),
        port: 1883,
    };
    let client = mqtt::get_client(config).await;
    return client;
}
