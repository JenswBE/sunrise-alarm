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
        self.config = settings.get()
        if not self.config.MOCK:
            self.strip = PixelStrip(
                num=self.config.LED_COUNT,
                pin=self.config.LED_GPIO_PIN,
                strip_type=self.config.LED_STRIP_TYPE
            )
            self.strip.begin()

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

    def update(self):
        """Set all leds to the current color"""
        # Skip if mocked
        if self.config.MOCK:
            return

        # Update leds
        color = _color_from_preset(self._color)
        for led_index in range(self.strip.numPixels()):
            self.strip.setPixelColorRGB(led_index, *color)
        self.strip.setBrightness(self._brightness)
        self.strip.show()

    def cleanup(self):
        """Cleanup on exit"""
        self.set_black()

    def set_color(self, color: PresetColor, brightness: int = 100):
        """Set all leds to a specific color and brightness (0 - 255)"""
        self._color = color
        self._brightness = brightness
        self.update()

    def set_black(self):
        """Turn off all leds"""
        self._color = PresetColor.BLACK
        self._brightness = 0
        self.update()

    def start_sunrise_simulation(self):
        """Start simulating a sunrise"""
        # Check if already in sunrise
        if self._sunrise:
            return

        # Set initial state
        self.set_color(PresetColor.RED, 1)
        self._sunrise = True

        # Set timer to update sunrise
        config = settings.get()
        loop = asyncio.get_event_loop()
        self._update_sunrise_event = loop.call_later(
            callback=self.update_sunrise_simulation,
            delay=config.LIGHT_INCREASE_DURATION.seconds / 100,
        )

    def update_sunrise_simulation(self):
        """Update the sunrise simulation"""
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
        loop = asyncio.get_event_loop()
        self._update_sunrise_event = loop.call_later(
            callback=self.update_sunrise_simulation,
            delay=config.LIGHT_INCREASE_DURATION.seconds / 100,
        )

    def stop_sunrise_simulation(self):
        """Stop simulating a sunrise"""
        # Check if we are in a sunrise
        if not self._sunrise:
            return

        # Cancel sunrise update if still runnning
        if self._update_sunrise_event is not None:
            self._update_sunrise_event.cancel()
            self._update_sunrise_event = None

        # Turn light off
        self.set_black()
