# Sunrise Alarm

DIY alarm clock
![Result](schematics/result.jpg)
![Scheme](schematics/scheme.jpg)

## Local development

Dependencies:

- NodeJS
- Golang
- ALSA dev library
  - Fedora: `sudo dnf install alsa-lib-devel`
- GCC
  - Fedora: `sudo dnf group install "C Development Tools and Libraries" "Development Tools"`

```bash
# Update vendor libraries
./src/update_vendors.sh

# Change to source dir
cd src

# Auto-restart Sunrise Alarm on changes to the code
go install github.com/cespare/reflex@latest
reflex -s -G alarms.json go run ./cmd/

# OR ignoring "next-ring-time" calls
reflex -s -G alarms.json go run ./cmd/ 2>&1 | grep -v "next-ring-time"
```

**NOTE:** GUI optimized for [Raspberry Pi Touch Display](https://www.raspberrypi.com/products/raspberry-pi-touch-display/) (800x480px)

## Deployment

**Known issues**

- Onscreen keyboard not working since switch to Wayland on Rpi OS 12 (OS only, OSK in alarm settings works)

```bash
# Enable SSH access
sudo raspi-config

# Give GPU 256MB of memory
sudo raspi-config nonint do_memory_split 256

# Enable I2C (0 = enabled, 1 = disabled)
sudo raspi-config nonint do_i2c 0

# Disable swap
# Based on https://raspberrypi.stackexchange.com/questions/84390/how-to-permanently-disable-swap-on-raspbian-stretch-lite
sudo dphys-swapfile swapoff
sudo dphys-swapfile uninstall
sudo systemctl disable dphys-swapfile
sudo update-rc.d dphys-swapfile remove
sudo apt purge dphys-swapfile

# Write all logs to RAM (will be exported anyway)
# Based on https://raspberrypi.stackexchange.com/questions/124605/stop-logs-from-writing-to-var-log
sudo tee -a /etc/fstab <<EOF
tmpfs /tmp tmpfs defaults,noatime,mode=1777 0 0
tmpfs /var/tmp tmpfs defaults,noatime,mode=1777 0 0
tmpfs /var/log tmpfs defaults,noatime,mode=0755 0 0
tmpfs /var/spool tmpfs defaults,noatime,mode=1777 0 0
EOF

# Export logs
# Based on https://www.elastic.co/downloads/beats/journalbeat
wget -qO - https://artifacts.elastic.co/downloads/beats/journalbeat/journalbeat-7.15.2-arm64.deb -O journalbeat.deb
sudo dpkg -i ./journalbeat.deb
rm ./journalbeat.deb
sudo nano /etc/journalbeat/journalbeat.yml
# 1. Comment key "output.elasticsearch" and children
# 2. Uncomment key "output.logstash" and configure child "hosts" to Graylog server
sudo systemctl enable --now journalbeat

# Upgrade system
sudo apt update
sudo apt dist-upgrade -y
sudo reboot # Needed for "Log to RAM" anyway

# Install dependencies
sudo apt install -y libasound2-dev

# Install Go
# Update below to latest version at https://go.dev/dl/
# See https://go.dev/doc/install for official instructions
GO_URL="https://go.dev/dl/go1.21.3.linux-arm64.tar.gz"
wget -O go.linux-arm64.tar.gz "${GO_URL:?}"
sudo rm -rf /usr/local/go ~/go || true # Removes old install
sudo tar -C /usr/local -xzf go.linux-arm64.tar.gz
rm go.linux-arm64.tar.gz
sudo tee /etc/profile.d/add-go-to-path.sh <<EOF
export PATH=\$PATH:/usr/local/go/bin
export GOPATH=~/go
EOF
source /etc/profile.d/add-go-to-path.sh

# Install NodeJS
# Based on https://github.com/nodesource/distributions#installation-instructions
NODE_VERSION=20
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
sudo tee /etc/apt/sources.list.d/nodesource.list <<EOF
deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_${NODE_VERSION:?}.x nodistro main
EOF
sudo apt-get update
sudo apt-get install nodejs -y

# Clone this repo
git config --global pull.ff only
git clone https://github.com/JenswBE/sunrise-alarm
sudo chown -R root:root sunrise-alarm
cd sunrise-alarm

# Setup Sunrise Alarm
# It's intended below variables are not escaped!
sudo mkdir -p /etc/systemd/system
sudo tee /etc/systemd/system/sunrise-alarm.service <<EOF
[Unit]
Description=Sunrise Alarm

[Service]
TimeoutSec=10min
WorkingDirectory=${HOME:?}/sunrise-alarm
Environment="GOPATH=${GOPATH:?}"
Environment="HOME=${HOME:?}"
Environment="PATH=${PATH:?}"
Environment="DEBUG=true"
ExecStartPre=-$(which git) config --global --add safe.directory ${HOME:?}/sunrise-alarm
ExecStartPre=-$(which git) pull
ExecStartPre=-$(which bash) src/update_vendors.sh
ExecStartPre=$(which bash) -c "cd src; go build -o ../sunrise-alarm ./cmd/"
ExecStart=${HOME:?}/sunrise-alarm/sunrise-alarm

[Install]
WantedBy=default.target
EOF
sudo systemctl daemon-reload
sudo systemctl enable --now sunrise-alarm

# Run Ansible
cd deployment
LC_ALL=C.UTF-8 ansible-galaxy collection install -r requirements.yml
LC_ALL=C.UTF-8 ansible-playbook main.yml

# Reboot
sudo reboot

# Open GUI
xdg-open "http://localhost:8123"
```
