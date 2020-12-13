use std::{env, path::PathBuf};

use url::Url;

#[derive(Debug, Clone)]
pub struct WarpConfig {
    pub port: u16,
}

impl WarpConfig {
    pub fn from_env(default_port: u16) -> Self {
        Self {
            port: parse_port("WARP_PORT", default_port),
        }
    }
}

#[derive(Debug, Clone)]
pub struct HostsConfig {
    pub srv_alarm: Url,
    pub srv_config: Url,
    pub srv_physical: Url,
    pub srv_audio: Url,
}

impl HostsConfig {
    pub fn from_env() -> Self {
        Self {
            srv_alarm: parse_url("HOST_SRV_ALARM", "http://localhost:8000"),
            srv_config: parse_url("HOST_SRV_CONFIG", "http://localhost:8001"),
            srv_physical: parse_url("HOST_SRV_PHYSICAL", "http://localhost:8002"),
            srv_audio: parse_url("HOST_SRV_AUDIO", "http://localhost:8003"),
        }
    }
}

#[derive(Debug, Clone)]
pub struct MqttConfig {
    pub host: String,
    pub port: u16,
}

impl MqttConfig {
    pub fn from_env() -> Self {
        Self {
            host: parse_string("MQTT_BROKER_HOST", "localhost"),
            port: parse_port("MQTT_BROKER_PORT", 1883),
        }
    }
}

// =============================================
// =                  HELPERS                  =
// =============================================
pub fn parse_string(env_var: &str, default: &str) -> String {
    env::var(env_var).unwrap_or(default.to_string())
}

pub fn parse_port(env_var: &str, default: u16) -> u16 {
    let port = env::var(env_var);
    if let Ok(port) = port {
        port.parse()
            .unwrap_or_else(|_| panic!("Provided {} is not a valid port: {}", env_var, port))
    } else {
        default
    }
}

pub fn parse_url(env_var: &str, default: &str) -> Url {
    let env_value = parse_string(env_var, default);
    Url::parse(&env_value)
        .unwrap_or_else(|_| panic!("Provided {} is not a valid URL: {}", env_var, env_value))
}

pub fn parse_path(env_var: &str, default: &str) -> PathBuf {
    let env_value = parse_string(env_var, default);
    env_value
        .parse()
        .unwrap_or_else(|_| panic!("Provided {} is not a valid path: {}", env_var, env_value))
}
