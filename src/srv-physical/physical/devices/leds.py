"""This module contains helpers to work with RGB leds"""

import asyncio
import atexit
import logging
from collections import namedtuple
from enum import Enum

from physical.helpers import settings

if not settings.get().MOCK:
    from rpi_ws281x import PixelStrip

Color = namedtuple('Color', 'red green blue')
BLACK = Color(0, 0, 0)
RED = Color(255, 0, 0)
ORANGE = Color(255, 100, 0)
YELLOW = Color(255, 255, 0)
WARM_WHITE = Color(239, 197, 59)


class PresetColor(str, Enum):
    BLACK = 'BLACK'
    RED = 'RED'
    ORANGE = 'ORANGE'
    YELLOW = 'YELLOW'
    WARM_WHITE = 'WARM_WHITE'


def _color_from_preset(preset: PresetColor) -> Color:
    return globals()[preset.value]


class Leds:
    """Helper class to work with RGB leds"""

    def __init__(self):
        # Create led strip
        self._is_real = not settings.get().MOCK
        self._strip = self._new_pixel_strip()
        self._loop = asyncio.get_event_loop()

        # Set initial color and brightness
        self._color = PresetColor.BLACK
        self._brightness = 0
        self.update()

        # Init sunrise
        self._sunrise = False
        self._update_sunrise_event = None

        # Register cleanup on exit
        atexit.register(self.cleanup)

    @property
    def color(self) -> PresetColor:
        return self._color

    @property
    def brightness(self) -> int:
        return self._brightness

    def _new_pixel_strip(self):
        if self._is_real:
            config = settings.get()
            strip = PixelStrip(
                num=config.LED_COUNT,
                pin=config.LED_GPIO_PIN,
                strip_type=config.LED_STRIP_TYPE
            )
            strip.begin()
            return strip

    def update(self):
        """Set all leds to the current color"""
        if self._is_real:
            color = _color_from_preset(self._color)
            for led_index in range(self._strip.numPixels()):
                self._strip.setPixelColorRGB(led_index, *color)
            self._strip.setBrightness(self._brightness)
            self._strip.show()

    def cleanup(self):
        """Cleanup on exit"""
        self.set_black()

    def set_color(self, color: PresetColor, brightness: int = 100):
        """Set all leds to a specific color and brightness (0 - 255)"""
        logging.info("Set leds to color %s with brightness %d",
                     color, brightness)
        self._color = color
        self._brightness = brightness
        self.update()

    def set_black(self):
        """Turn off all leds"""
        logging.info("Turn leds off")
        self._color = PresetColor.BLACK
        self._brightness = 0
        self.update()

    def start_sunrise_simulation(self):
        """Start simulating a sunrise"""
        logging.info("Start sunrise simulation")

        # Check if already in sunrise
        if self._sunrise:
            logging.info("We are already in a sunrise, ignore start sunrise")
            return

        # Set initial state
        self.set_color(PresetColor.RED, 1)
        self._sunrise = True

        # Set timer to update sunrise
        config = settings.get()
        self._update_sunrise_event = self._loop.call_later(
            callback=self.update_sunrise_simulation,
            delay=config.LIGHT_INCREASE_DURATION.seconds / 100,
        )

    def update_sunrise_simulation(self):
        """Update the sunrise simulation"""
        logging.info("Update sunrise simulation")

        # Derive color from brightness
        brightness = self._brightness + 1
        if brightness > 90:
            color = PresetColor.WARM_WHITE
        elif brightness > 60:
            color = PresetColor.YELLOW
        elif brightness > 30:
            color = PresetColor.ORANGE
        else:
            color = PresetColor.RED

        # Update leds
        self.set_color(color, brightness)

        # Cancel updating sunrise at brightness 100
        if brightness >= 100 and self._update_sunrise_event is not None:
            self._update_sunrise_event.cancel()
            self._update_sunrise_event = None

        # Reschedule
        config = settings.get()
        self._update_sunrise_event = self._loop.call_later(
            callback=self.update_sunrise_simulation,
            delay=config.LIGHT_INCREASE_DURATION.seconds / 100,
        )

    def stop_sunrise_simulation(self):
        """Stop simulating a sunrise"""
        logging.info("Stop sunrise simulation")

        # Check if we are in a sunrise
        if not self._sunrise:
            logging.info("We are not in a sunrise, ignore stop sunrise")
            return

        # Cancel sunrise update if still runnning
        if self._update_sunrise_event is not None:
            self._update_sunrise_event.cancel()
            self._update_sunrise_event = None

        # Reset state
        self.set_black()
        self._sunrise = False
