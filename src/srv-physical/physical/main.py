"""Main entry point for srv-physical"""

import logging

from fastapi import FastAPI
from starlette.responses import RedirectResponse
from starlette.middleware.cors import CORSMiddleware

from physical.helpers import (mqtt, settings)
from physical.rest import (
    mock as router_mock,
)

# Start app
app = FastAPI(
    title='srv-physical',
    description='Service to handle buttons, leds, display backlight, ...',
)

# Add routers
# Mock
app.include_router(
    router_mock.router,
    prefix='/mock',
    tags=['Mock']
)


@app.on_event('startup')
async def setup_service():
    """Setup srv-physical"""
    logging.info("MQTT client: Setting up ...")
    app.state.mqtt = await mqtt.get()
    logging.info("MQTT client: Setup finished")

    config = settings.get()
    if not config.MOCK:
        from physical.devices import (button, display, leds)
        app.state.devices = {
            "button": button.Button(config.BUTTON_GPIO_PIN),
            "display": display.Display(),
            "leds": leds.Leds(),
        }


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
