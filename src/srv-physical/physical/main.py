"""Main entry point for srv-physical"""

import logging
from physical import devices

from fastapi import FastAPI
from starlette.responses import RedirectResponse
from starlette.middleware.cors import CORSMiddleware

from physical.devices import button, buzzer, display, leds, temp_humid
from physical.devices.common import Devices
from physical.helpers import mqtt, settings
from physical.rest import (
    backlight as router_backlight,
    buzzer as router_buzzer,
    leds as router_leds,
    mock as router_mock,
)

# Start app
app = FastAPI(
    title='srv-physical',
    description='Service to handle buttons, leds, display backlight, ...',
)

# Add routers
# Backlight
app.include_router(
    router_backlight.router,
    prefix='/backlight',
    tags=['Backlight']
)

# Buzzer
app.include_router(
    router_buzzer.router,
    prefix='/buzzer',
    tags=['Buzzer']
)

# Leds
app.include_router(
    router_leds.router,
    prefix='/leds',
    tags=['LEDs']
)

# Mock
app.include_router(
    router_mock.router,
    prefix='/mock',
    tags=['Mock']
)


@app.on_event('startup')
async def setup_service():
    """Setup srv-physical"""
    # Setup MQTT
    logging.info("MQTT client: Setting up ...")
    app.state.mqtt = await mqtt.get()
    logging.info("MQTT client: Setup finished")

    # Setup devices
    config = settings.get()
    top_button = button.Button(config.BUTTON_GPIO_PIN)
    top_button.set_short_press_callback(
        app.state.mqtt.publish_button_pressed)
    top_button.set_long_press_callback(
        app.state.mqtt.publish_button_long_pressed)
    th_sensor = temp_humid.TempHumid(config.TEMP_HUMID_GPIO_PIN)
    th_sensor.set_new_reading_callback(
        app.state.mqtt.publish_temp_humid_updated)
    app.state.devices = Devices(
        button=top_button,
        buzzer=buzzer.Buzzer(config.BUZZER_GPIO_PIN),
        display=display.Display(),
        leds=leds.Leds(),
        th_sensor=th_sensor,
    )

    # Use DPMS instead of custom sleep for time being
    app.state.devices.display.enable_keep_awake()


@app.on_event("shutdown")
async def shutdown_service():
    """Shutdown srv-physical"""
    logging.info("MQTT client: Closing ...")
    await app.state.mqtt.stop()
    logging.info("MQTT client: Closed")


# Add CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get('/', include_in_schema=False)
async def redirect_to_docs():
    """Redirect root page to docs."""
    return RedirectResponse(url='/docs')
