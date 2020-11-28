use rumqttc::{AsyncClient, Event, MqttOptions, Packet, Publish};

use crate::models::Context;
use crate::time;
use sunrise_common::mqtt;

const CLIENT_ID: &str = "srv-alarm";

pub async fn get_client(ctx: Context) -> AsyncClient {
    // Build client
    let mut options = MqttOptions::new(CLIENT_ID, &ctx.config.mqtt.host, ctx.config.mqtt.port);
    options.set_keep_alive(5);
    let (mqtt_client, mut eventloop) = AsyncClient::new(options, 10);

    // Start client loop
    tokio::task::spawn(async move {
        loop {
            let notification = eventloop.poll().await.unwrap();
            notification_handler(ctx.clone(), notification).await;
        }
    });

    // Subscribe to topics
    mqtt::subscribe(&mqtt_client, mqtt::TOPIC_ALARMS_CHANGED).await;

    // Create client successful
    return mqtt_client;
}

async fn notification_handler(ctx: Context, notification: Event) {
    // Debug logging
    log::debug!("MQTT notification received: {:?}", notification);

    // Only handle incoming publish notifications
    if let Event::Incoming(event) = notification {
        if let Packet::Publish(packet) = event {
            match packet.topic.as_str() {
                mqtt::TOPIC_ALARMS_CHANGED => handle_alarms_changed(ctx, packet).await,
                _ => log::error!("Unhandled MQTT topic: {}", packet.topic),
            }
        }
    }
}

async fn handle_alarms_changed(ctx: Context, packet: Publish) {
    let alarms = mqtt::parse_alarms_changed(packet).alarms;
    ctx.set_alarms(alarms);
    time::update_next_alarms(ctx);
}
