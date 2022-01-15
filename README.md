# Sunrise Alarm

[![GitHub Repo](https://img.shields.io/badge/GitHub-repo-brightgreen?logo=github)](https://github.com/JenswBE/sunrise-alarm)

DIY alarm clock using microservices
![Result](schematics/result.jpg)
![Scheme](schematics/scheme.jpg)

## Services

| Service       | Description                                         |                                                                      Links                                                                      |  Dev port  |  Language  |   Frameworks    |
| ------------- | --------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------: | :--------: | :--------: | :-------------: |
| srv-alarm     | Main logic of the alarm                             |   [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/r/jenswbe/sunrise-alarm-srv-alarm)   |    8000    |    Rust    |      Warp       |
| srv-config    | Configuration management                            |  [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/r/jenswbe/sunrise-alarm-srv-config)   |    8001    |    Rust    | Warp, Rustbreak |
| srv-physical  | Interacts with physical features: button, leds, ... | [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/r/jenswbe/sunrise-alarm-srv-physical)  |    8002    |     Go     |       Gin       |
| srv-audio     | Alarm sound handling                                |   [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/r/jenswbe/sunrise-alarm-srv-audio)   |    8003    |    Rust    |   Warp, Rodio   |
| api-watchface | REST API for watchface UI                           |              [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/r/_/traefik)              |    8004    |    N/A     |     Traefik     |
| gui-watchface | Web UI for touchscreen                              | [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/r/jenswbe/sunrise-alarm-gui-watchface) |    8080    | Typescript |     Nuxt.js     |
| mosquitto     | MQTT broker                                         |          [![DockerHub Repo](https://img.shields.io/badge/DockerHub-repo-blue?logo=docker)](https://hub.docker.com/_/eclipse-mosquitto)          | 1883, 9001 |    N/A     |    Mosquitto    |

## Development

Start Docker Compose with following command:

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose build
docker-compose up -d
```

## Setup

```bash
# Enable SSH access
sudo raspi-config

# Give GPU 256MB of memory
sudo raspi-config nonint do_memory_split 256

# Enable I2C (0 = enabled, 1 = disabled)
sudo raspi-config nonint do_i2c 0

# Set hostname
sudo hostnamectl set-hostname sunrise

# Disable swap
# Based on https://raspberrypi.stackexchange.com/questions/84390/how-to-permanently-disable-swap-on-raspbian-stretch-lite
sudo dphys-swapfile swapoff
sudo dphys-swapfile uninstall
sudo systemctl disable dphys-swapfile
sudo update-rc.d dphys-swapfile remove
sudo apt purge dphys-swapfile

# Write all logs to RAM (will be exported with Promtail anyway)
# Based on https://raspberrypi.stackexchange.com/questions/124605/stop-logs-from-writing-to-var-log
sudo tee -a /etc/fstab <<EOF
tmpfs /tmp tmpfs defaults,noatime,mode=1777 0 0
tmpfs /var/tmp tmpfs defaults,noatime,mode=1777 0 0
tmpfs /var/log tmpfs defaults,noatime,mode=0755 0 0
tmpfs /var/spool tmpfs defaults,noatime,mode=1777 0 0
EOF

# Upgrade system
sudo apt update
sudo apt dist-upgrade -y
sudo reboot # Needed for "Log to RAM" anyway

# Install dependencies
sudo apt install -y firefox-esr onboard

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker pi

# Install Docker Compose
DOCKER_COMPOSE_VERSION=v2.1.1 # See https://github.com/docker/compose/releases for latest version
DOCKER_COMPOSE_URL="https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-linux-armv7"
sudo mkdir -p /usr/local/lib/docker/cli-plugins
sudo curl -SL ${DOCKER_COMPOSE_URL} -o /usr/local/lib/docker/cli-plugins/docker-compose
sudo chmod +x /usr/local/lib/docker/cli-plugins/docker-compose

# Clone this repo
git clone https://github.com/JenswBE/sunrise-alarm
cd sunrise-alarm

# Configure rsyslog
sudo cp deployment/rsyslog-40-docker.conf /etc/rsyslog.d/40-docker.conf
sudo systemctl restart rsyslog
sudo cp deployment/logrotate-docker /etc/logrotate.d/docker

# Configure docker to log to syslog
sudo cp deployment/docker-daemon.json /etc/docker/daemon.json
sudo systemctl restart docker

# Configure promtail: Replace placeholders <LOKI_*>
cp deployment/promtail-config.yml.template deployment/promtail-config.yml
nano deployment/promtail-config.yml

# Configure screen timeout
mkdir -p ~/.config/autostart
tee ~/.config/autostart/screen-timeout.desktop <<EOF
[Desktop Entry]
Type=Application
Name=Set screen timeout
Exec=/usr/bin/xset dpms 60 60 60
EOF

# Configure PulseAudio
# Based on:
#   - http://manpages.ubuntu.com/manpages/xenial/man5/default.pa.5.html
#   - https://wiki.archlinux.org/index.php/PulseAudio/Examples#Allowing_multiple_users_to_use_PulseAudio_at_the_same_time
tee ~/.config/pulse/default.pa <<EOF
#!/usr/bin/pulseaudio -nF

.include /etc/pulse/default.pa

# Sunrise alarm
load-module module-native-protocol-unix auth-anonymous=1 socket=/tmp/pa-sunrise-alarm.socket
EOF

# Reboot
sudo reboot

# Configure alarm
cp deployment/.env.template deployment/.env
nano deployment/.env

# Start alarm
docker compose up -d
xdg-open "http://localhost:8080"
```

## Firefox Add-ons:

### ScrollAnywhere

- General: Enable `Left button`
- General: Enable `Scroll on text` + `Double click`
- General: Set scroll style `Grab and Drag`
- Scrollbars: Set look `Thin`
- Scrollbars: Set background to black and slider to dark grey
