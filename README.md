# Sunrise Alarm
[![GitHub Repo](https://img.shields.io/badge/GitHub-repo-brightgreen?logo=github)](https://github.com/JenswBE/sunrise-alarm)

DIY alarm clock using microservices
![Result](schematics/result.jpg)
![Scheme](schematics/scheme.jpg)

## Services
| Service       | Description                                         | Links                                                                                                                                           | Dev port   | Language   | Frameworks      |
|---------------|-----------------------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------:|:----------:|:----------:|:---------------:|
| srv-alarm     | Main logic of the alarm                             |                                                                                                                                                 | 8000       | Rust       | Warp            |
| srv-config    | Configuration management                            | [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/r/jenswbe/sunrise-alarm-srv-config)    | 8001       | Rust       | Warp, Rustbreak |
| srv-physical  | Interacts with physical features: button, leds, ... | [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/r/jenswbe/sunrise-alarm-srv-physical)  | 8002       | Python     | FastAPI         |
| srv-audio     | Alarm sound handling                                | [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/r/jenswbe/sunrise-alarm-srv-audio)     | 8003       | Rust       | Warp, Rodio     |
| api-watchface | REST API for watchface UI                           | [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/r/devopsfaith/krakend)                 | 8004       | N/A        | KrakenD         |
| gui-watchface | Web UI for touchscreen                              | [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/r/jenswbe/sunrise-alarm-gui-watchface) | 8080       | Javascript | Vue.js          |
| mosquitto     | MQTT broker                                         | [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/_/eclipse-mosquitto)                   | 1883, 9001 | N/A        | N/A             |

## Setup

### PulseAudio
Based on [Arch wiki](https://wiki.archlinux.org/index.php/PulseAudio/Examples#Allowing_multiple_users_to_use_PulseAudio_at_the_same_time)
and [Ubuntu man pages](http://manpages.ubuntu.com/manpages/xenial/man5/default.pa.5.html)
```bash
tee ~/.config/pulse/default.pa <<EOF
#!/usr/bin/pulseaudio -nF

.include /etc/pulse/default.pa

# Sunrise alarm
load-module module-native-protocol-unix auth-anonymous=1 socket=/tmp/pa-sunrise-alarm.socket
EOF
```
