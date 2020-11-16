# Sunrise Alarm
[![GitHub Repo](https://img.shields.io/badge/GitHub-repo-brightgreen?logo=github)](https://github.com/JenswBE/sunrise-alarm)

DIY alarm clock using microservices
![Result](schematics/result.jpg)
![Scheme](schematics/scheme.jpg)

## Services
| Service       | Description                                         | Links                                                                                                                                        | Dev port | Language   | Frameworks      |
|---------------|-----------------------------------------------------|:--------------------------------------------------------------------------------------------------------------------------------------------:|:--------:|:----------:|:---------------:|
| srv-alarm     | Main logic of the alarm                             |                                                                                                                                              | 8000     | Rust       | Warp            |
| srv-config    | Configuration management                            | [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/r/jenswbe/sunrise-alarm-srv-config) | 8001     | Rust       | Warp, Rustbreak |
| srv-physical  | Interacts with physical features: button, leds, ... |                                                                                                                                              | 8002     | Python     | FastAPI         |
| srv-audio     | Alarm sound handling                                | [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/r/jenswbe/sunrise-alarm-srv-audio)  | 8003     | Rust       | Warp, Rodio     |
| api-watchface | REST API for watchface UI                           |                                                                                                                                              | 8004     | Rust       | Warp            |
| gui-watchface | Web UI for touchscreen                              |                                                                                                                                              | 8080     | Javascript | Vue.js          |
| mosquitto     | MQTT broker                                         | [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/_/eclipse-mosquitto)                | 1883     | N/A        | N/A             |

## Setup

### PulseAudio
Based on [Arch wiki](https://wiki.archlinux.org/index.php/PulseAudio/Examples#Allowing_multiple_users_to_use_PulseAudio_at_the_same_time)
and [Ubuntu man pages](http://manpages.ubuntu.com/manpages/xenial/man5/default.pa.5.html)
```
cp /etc/pulse/default.pa ~/.config/pulse/default.pa
tee --append ~/.config/pulse/default.pa <<EOF

# Sunrise alarm
load-module module-native-protocol-unix auth-anonymous=1 socket=/tmp/pa-sunrise-alarm.socket
EOF
```

## Cross-compilation
https://users.rust-lang.org/t/static-cross-build-for-arm/9100/2

```bash
sudo apt-get install -qq gcc-arm-linux-gnueabihf
rustup target add armv7-unknown-linux-musleabihf
cd src/srv-config
export CARGO_TARGET_ARMV7_UNKNOWN_LINUX_MUSLEABIHF_LINKER=arm-linux-gnueabihf-gcc
export CC_armv7_unknown_linux_musleabihf=arm-linux-gnueabihf-gcc
cargo build --target armv7-unknown-linux-musleabihf
```