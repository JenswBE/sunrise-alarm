use std::path::PathBuf;
use std::sync::{Arc, Mutex};

use sunrise_common::config::{parse_path, MqttConfig, WarpConfig};

pub type Music = Arc<Mutex<rodio::Sink>>;

#[derive(Debug, Clone)]
pub struct Config {
    pub music_dir: PathBuf,
    pub warp: WarpConfig,
    pub mqtt: MqttConfig,
}

impl Config {
    pub fn from_env() -> Self {
        Self {
            music_dir: parse_path("MUSIC_DIR_PATH", "../../data/music"),
            warp: WarpConfig::from_env(8003),
            mqtt: MqttConfig::from_env(),
        }
    }
}
