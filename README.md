# Sunrise Alarm
DIY alarm clock using microservices
![Result](schematics/result.jpg)
![Scheme](schematics/scheme.jpg)

## Services
| Service       | Description                                         | Dev port | Language   | Frameworks      |
|---------------|-----------------------------------------------------|:--------:|:----------:|:---------------:|
| srv-alarm     | Main logic of the alarm                             | 8000     | Rust       | Warp            |
| srv-config    | Configuration management                            | 8001     | Rust       | Warp, Rustbreak |
| srv-physical  | Interacts with physical features: button, leds, ... | 8002     | Python     | FastAPI         |
| srv-audio     | Alarm sound handling                                | 8003     | Rust       | Warp, SFML      |
| api-watchface | REST API for watchface UI                           | 8004     | Rust       | Warp            |
| gui-watchface | Web UI for touchscreen                              | 8080     | Javascript | Vue.js          |
| mosquitto     | MQTT broker                                         | 1883     | N/A        | N/A             |

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