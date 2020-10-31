# Sunrise Alarm
DIY alarm clock using microservices
![Result](schematics/result.jpg)
![Scheme](schematics/scheme.jpg)

## Services
- srv-config: Configuration management (Rust)
- srv-alarm: Contains main logic of the alarm (Rust)
- srv-physical: Service which interacts with physical alarm features: buttons, LED's, sound, ... (Python FastAPI)
- gui-watchface: Nuxt.js application which serves as GUI for touchscreen
- api-watchface: REST API to support watchface GUI (Rust)