use rumqttc::{AsyncClient, Event, MqttOptions, Packet, Publish};

use crate::models::{Config, State};
use crate::time;
use sunrise_common::mqtt;

const CLIENT_ID: &str = "srv-alarm";

pub async fn get_client(config: &Config, state: State) -> AsyncClient {
    // Build client
    let mut options = MqttOptions::new(CLIENT_ID, config.mqtt.host.clone(), config.mqtt.port);
    options.set_keep_alive(5);
    let (mqtt_client, mut eventloop) = AsyncClient::new(options, 10);

    // Start client loop
    let loop_config = config.clone();
    tokio::task::spawn(async move {
        loop {
            let notification = eventloop.poll().await.unwrap();
            notification_handler(notification, state.clone(), loop_config.clone()).await;
        }
    });

    // Subscribe to topics
    mqtt::subscribe(&mqtt_client, mqtt::TOPIC_ALARMS_CHANGED).await;

    // Create client successful
    return mqtt_client;
}

async fn notification_handler(notification: Event, state: State, config: Config) {
    // Debug logging
    log::debug!("MQTT notification received: {:?}", notification);

    // Only handle incoming publish notifications
    if let Event::Incoming(event) = notification {
        if let Packet::Publish(packet) = event {
            match packet.topic.as_str() {
                mqtt::TOPIC_ALARMS_CHANGED => handle_alarms_changed(packet, state, config).await,
                _ => log::error!("Unhandled MQTT topic: {}", packet.topic),
            }
        }
    }
}

async fn handle_alarms_changed(packet: Publish, state: State, config: Config) {
    {
        let mut state = state.lock().unwrap();
        state.alarms = mqtt::parse_alarms_changed(packet).alarms;
    }
    time::update_next_alarms(state, &config.alarm)
}
