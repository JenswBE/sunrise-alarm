use sunrise_common::alarm::Alarm;

use crate::models::Context;

pub async fn get_alarms(ctx: Context) -> Vec<Alarm> {
    let url = ctx.config.hosts.srv_config.join("alarms").unwrap();
    ctx.client
        .get(url)
        .send()
        .await
        .unwrap()
        .json()
        .await
        .unwrap()
}

pub async fn update_alarm(ctx: Context, alarm: Alarm) {
    let path = format!("alarms/{}", alarm.id);
    let url = ctx.config.hosts.srv_config.join(&path).unwrap();
    ctx.client.put(url).json(&alarm).send().await.unwrap();
}
