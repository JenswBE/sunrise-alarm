use std::path::PathBuf;
use std::sync::{Arc, Mutex};

use sunrise_common::mqtt::MqttConfig;

pub type Music = Arc<Mutex<rodio::Sink>>;

#[derive(Debug, Clone)]
pub struct Config {
    pub music_dir: PathBuf,
    pub warp_port: u16,
    pub mqtt_config: MqttConfig,
}
