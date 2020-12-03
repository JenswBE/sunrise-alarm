use reqwest::{Error, Response, Url};
use serde::{Deserialize, Serialize};

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
        .map_err(format_error("Failed to GET alarms", url.clone()))?;

    error_for_status(response, "Failed to GET alarms")
        .await?
        .json()
        .await
        .map_err(format_error("Failed to parse alarms after GET", url))
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
        .map(|_| ()) // We don't care about response body
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
    pub fn black() -> Leds {
        Leds {
            color: LedsColor::Black,
            brightness: None,
        }
    }

    pub fn night_light() -> Leds {
        Leds {
            color: LedsColor::WarmWhite,
            brightness: None,
        }
    }

    pub fn night_light_dark() -> Leds {
        Leds {
            color: LedsColor::Orange,
            brightness: Some(2),
        }
    }
}

impl Default for Leds {
    fn default() -> Self {
        Leds {
            color: LedsColor::Black,
            brightness: None,
        }
    }
}

#[serde(rename_all = "SCREAMING_SNAKE_CASE")]
#[derive(Debug, PartialEq, Serialize, Deserialize)]
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
        .map(|_| ()) // We don't care about response body
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
        .map(|_| ()) // We don't care about response body
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
