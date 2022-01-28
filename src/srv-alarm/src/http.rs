use reqwest::{Error, Response, Url};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

use crate::models::Context;
use sunrise_common::alarm::Alarm;

// ==============================================
// =                   ALARMS                   =
// ==============================================

/// srv-config: GET /alarms
pub async fn get_alarms(ctx: &Context) -> Result<Vec<Alarm>, String> {
    let url = ctx.config.hosts.srv_config.join("alarms").unwrap();
    let response = ctx
        .client
        .get(url.clone())
        .send()
        .await
        .map_err(format_error("Failed to GET alarms list", url.clone()))?;

    error_for_status(response, "Failed to GET alarms list")
        .await?
        .json()
        .await
        .map_err(format_error("Failed to parse alarms after GET list", url))
}

/// srv-config: GET /alarms/{alarm_id}
pub async fn get_alarm(ctx: &Context, alarm_id: Uuid) -> Result<Alarm, String> {
    let path = format!("alarms/{}", alarm_id);
    let url = ctx.config.hosts.srv_config.join(&path).unwrap();
    let response = ctx
        .client
        .get(url.clone())
        .send()
        .await
        .map_err(format_error("Failed to GET single alarm", url.clone()))?;

    error_for_status(response, "Failed to GET single alarm")
        .await?
        .json()
        .await
        .map_err(format_error("Failed to parse alarm after GET single", url))
}

/// srv-config: PUT /alarms/{alarm_id}
pub async fn update_alarm(ctx: &Context, alarm: Alarm) -> Result<(), String> {
    let path = format!("alarms/{}", alarm.id);
    let url = ctx.config.hosts.srv_config.join(&path).unwrap();
    let response = ctx
        .client
        .put(url.clone())
        .json(&alarm)
        .send()
        .await
        .map_err(format_error("Failed to PUT alarm", url))?;

    error_for_status(response, "Failed to PUT alarm")
        .await
        .map(drop)
}

// ==============================================
// =                  BACKLIGHT                 =
// ==============================================

/// srv-physical: PUT /backlight/lock
pub async fn lock_backlight(ctx: &Context) -> Result<(), String> {
    let url = ctx
        .config
        .hosts
        .srv_physical
        .join("backlight/lock")
        .unwrap();
    let response = ctx
        .client
        .put(url.clone())
        .send()
        .await
        .map_err(format_error("Failed to PUT backlight/lock", url))?;

    error_for_status(response, "Failed to PUT backlight/lock")
        .await
        .map(drop)
}

/// srv-physical: DELETE /backlight/lock
pub async fn unlock_backlight(ctx: &Context) -> Result<(), String> {
    let url = ctx
        .config
        .hosts
        .srv_physical
        .join("backlight/lock")
        .unwrap();
    let response = ctx
        .client
        .delete(url.clone())
        .send()
        .await
        .map_err(format_error("Failed to DELETE backlight/lock", url))?;

    error_for_status(response, "Failed to DELETE backlight/lock")
        .await
        .map(drop)
}

// ==============================================
// =                   BUZZER                   =
// ==============================================

/// srv-physical: PUT /buzzer
pub async fn start_buzzer(ctx: &Context) -> Result<(), String> {
    let url = ctx.config.hosts.srv_physical.join("buzzer").unwrap();
    let response = ctx
        .client
        .put(url.clone())
        .send()
        .await
        .map_err(format_error("Failed to PUT buzzer", url))?;

    error_for_status(response, "Failed to PUT buzzer")
        .await
        .map(drop)
}

/// srv-physical: DELETE /buzzer
pub async fn stop_buzzer(ctx: &Context) -> Result<(), String> {
    let url = ctx.config.hosts.srv_physical.join("buzzer").unwrap();
    let response = ctx
        .client
        .delete(url.clone())
        .send()
        .await
        .map_err(format_error("Failed to DELETE buzzer", url))?;

    error_for_status(response, "Failed to DELETE buzzer")
        .await
        .map(drop)
}

// ==============================================
// =                    MUSIC                   =
// ==============================================

/// srv-audio: POST /music
pub async fn start_music(ctx: &Context) -> Result<(), String> {
    let url = ctx.config.hosts.srv_audio.join("music").unwrap();
    let response = ctx
        .client
        .post(url.clone())
        .send()
        .await
        .map_err(format_error("Failed to POST music", url))?;

    error_for_status(response, "Failed to POST music")
        .await
        .map(drop)
}

/// srv-audio: DELETE /music
pub async fn stop_music(ctx: &Context) -> Result<(), String> {
    let url = ctx.config.hosts.srv_audio.join("music").unwrap();
    let response = ctx
        .client
        .delete(url.clone())
        .send()
        .await
        .map_err(format_error("Failed to DELETE music", url))?;

    error_for_status(response, "Failed to DELETE music")
        .await
        .map(drop)
}

/// srv-audio: POST /volume/increase
pub async fn increase_music_volume(ctx: &Context) -> Result<(), String> {
    let url = ctx.config.hosts.srv_audio.join("volume/increase").unwrap();
    let response = ctx
        .client
        .post(url.clone())
        .send()
        .await
        .map_err(format_error("Failed to POST volume/increase", url))?;

    error_for_status(response, "Failed to POST volume/increase")
        .await
        .map(drop)
}

// ==============================================
// =                    LEDS                    =
// ==============================================

#[derive(Debug, Serialize, Deserialize)]
pub struct Leds {
    color: LedsColor,

    #[serde(skip_serializing_if = "Option::is_none")]
    brightness: Option<u8>,
}

impl Leds {
    pub fn is_off(&self) -> bool {
        self.color == LedsColor::Black || self.brightness.unwrap_or(0) == 0
    }

    pub fn is_on(&self) -> bool {
        !self.is_off()
    }

    // Presets
    pub fn black() -> Self {
        Self {
            color: LedsColor::Black,
            brightness: None,
        }
    }

    pub fn night_light() -> Self {
        Self {
            color: LedsColor::WarmWhite,
            brightness: None,
        }
    }

    pub fn night_light_dark() -> Self {
        Self {
            color: LedsColor::Orange,
            brightness: Some(10),
        }
    }
}

impl Default for Leds {
    fn default() -> Self {
        Self {
            color: LedsColor::Black,
            brightness: None,
        }
    }
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "SCREAMING_SNAKE_CASE")]
pub enum LedsColor {
    Black,
    Red,
    Orange,
    Yellow,
    WarmWhite,
}

/// srv-physical: GET /leds
pub async fn get_leds(ctx: &Context) -> Result<Leds, String> {
    let url = ctx.config.hosts.srv_physical.join("leds").unwrap();
    let response = ctx
        .client
        .get(url.clone())
        .send()
        .await
        .map_err(format_error("Failed to GET leds", url.clone()))?;

    error_for_status(response, "Failed to GET leds")
        .await?
        .json()
        .await
        .map_err(format_error("Failed to parse leds after GET", url))
}

/// srv-physical: PUT /leds
pub async fn set_leds(ctx: &Context, req: &Leds) -> Result<(), String> {
    let url = ctx.config.hosts.srv_physical.join("leds").unwrap();
    let response = ctx
        .client
        .put(url.clone())
        .json(req)
        .send()
        .await
        .map_err(format_error("Failed to PUT leds", url))?;

    error_for_status(response, "Failed to PUT leds")
        .await
        .map(drop)
}

/// srv-physical: DELETE /leds
pub async fn set_leds_off(ctx: &Context) -> Result<(), String> {
    let url = ctx.config.hosts.srv_physical.join("leds").unwrap();
    let response = ctx
        .client
        .delete(url.clone())
        .send()
        .await
        .map_err(format_error("Failed to DELETE leds", url))?;

    error_for_status(response, "Failed to DELETE leds")
        .await
        .map(drop)
}

/// srv-physical: PUT /leds/sunrise
pub async fn start_sunrise(ctx: &Context) -> Result<(), String> {
    let url = ctx.config.hosts.srv_physical.join("leds/sunrise").unwrap();
    let response = ctx
        .client
        .put(url.clone())
        .send()
        .await
        .map_err(format_error("Failed to PUT leds/sunrise", url))?;

    error_for_status(response, "Failed to PUT leds/sunrise")
        .await
        .map(drop)
}

/// srv-physical: DELETE leds/sunrise
pub async fn stop_sunrise(ctx: &Context) -> Result<(), String> {
    let url = ctx.config.hosts.srv_physical.join("leds/sunrise").unwrap();
    let response = ctx
        .client
        .delete(url.clone())
        .send()
        .await
        .map_err(format_error("Failed to DELETE leds/sunrise", url))?;

    error_for_status(response, "Failed to DELETE leds/sunrise")
        .await
        .map(drop)
}

// ==============================================
// =                  HELPERS                   =
// ==============================================
fn format_error(message: &'static str, url: Url) -> Box<dyn Fn(Error) -> String> {
    Box::new(move |error: Error| {
        let msg = format!("{}: {} - {}", message, url, error);
        log::error!("{}", msg);
        msg
    })
}

async fn error_for_status(response: Response, message: &'static str) -> Result<Response, String> {
    let status = response.status();
    let url = response.url().clone();
    if status.is_client_error() || status.is_server_error() {
        let err = format!(
            "{}: {} {} - {:?}",
            message,
            url,
            status,
            response.text().await
        );
        log::error!("{}", err);
        return Err(err);
    }
    Ok(response)
}
