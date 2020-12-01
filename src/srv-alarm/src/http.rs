use futures_util::TryFutureExt;
use reqwest::Error;
use serde::{Deserialize, Serialize};

use crate::models::Context;
use sunrise_common::alarm::Alarm;

// ==============================================
// =                   ALARMS                   =
// ==============================================

/// srv-config: GET /alarms
pub async fn get_alarms(ctx: &Context) -> Result<Vec<Alarm>, Error> {
    let url = ctx.config.hosts.srv_config.join("alarms").unwrap();
    let response = ctx.client.get(url).send().await;
    if let Ok(response) = response {
        response
            .error_for_status()
            .map_err(|e| {
                log::error!("Failed to fetch alarms: {}", e);
                return e;
            })
            .ok();
    }

    response.json().await.map_err(|e| {
        log::error!("Failed to parse alarms: {}", e);
        return e;
    })
}

/// srv-config: PUT /alarms/{alarm_id}
pub async fn update_alarm(ctx: &Context, alarm: Alarm) -> Result<(), Error> {
    let path = format!("alarms/{}", alarm.id);
    let url = ctx.config.hosts.srv_config.join(&path).unwrap();
    ctx.client
        .put(url)
        .json(&alarm)
        .send()
        .await
        .and_then(|r| r.error_for_status())
        .map_err(|e| {
            log::error!("Failed to fetch alarms: {}", e);
            return e;
        })
        .map(|_| ()) // We don't care about response
}

// ==============================================
// =                    LEDS                    =
// ==============================================

#[derive(Debug, Serialize, Deserialize)]
pub struct Leds {
    color: LedsColor,
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
pub async fn get_leds(ctx: &Context) -> Result<Leds, Error> {
    let url = ctx.config.hosts.srv_physical.join("leds").unwrap();
    ctx.client
        .get(url)
        .send()
        .and_then(|r| r.json())
        .await
        .map_err(|e| {
            log::error!("Failed to get leds: {}", e);
            return e;
        })
}

/// srv-physical: PUT /leds
pub async fn set_leds(ctx: &Context, req: &Leds) -> Result<(), Error> {
    let url = ctx.config.hosts.srv_physical.join("leds").unwrap();
    ctx.client
        .put(url)
        .json(req)
        .send()
        .await
        .and_then(|r| r.error_for_status())
        .map_err(|e| {
            log::error!("Failed to set leds: {:?}", e);
            return e;
        })
        .map(|_| ()) // We don't care about response
}

/// srv-physical: DELETE /leds
pub async fn set_leds_off(ctx: &Context) -> Result<(), Error> {
    let url = ctx.config.hosts.srv_physical.join("leds").unwrap();
    ctx.client
        .delete(url)
        .send()
        .await
        .and_then(|r| r.error_for_status())
        .map_err(|e| {
            log::error!("Failed to set leds to off: {}", e);
            return e;
        })
        .map(|_| ()) // We don't care about response
}
