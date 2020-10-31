use srv_config::Alarm;
use uuid::Uuid;

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
