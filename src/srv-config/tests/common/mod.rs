use std::collections::HashMap;
use std::sync::Arc;

use rustbreak::PathDatabase;
use tempfile::NamedTempFile;
use uuid::Uuid;

use srv_config::database;
use srv_config::models;
use srv_config::Alarm;

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
