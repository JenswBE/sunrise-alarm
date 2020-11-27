"""This module contains helpers to work with RGB leds"""

import asyncio
import atexit
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


class PresetColorEnum(str, Enum):
    BLACK = 'BLACK'
    RED = 'RED'
    ORANGE = 'ORANGE'
    YELLOW = 'YELLOW'
    WARM_WHITE = 'WARM_WHITE'


class Leds:
    """Helper class to work with RGB leds"""

    def __init__(self):
        # Create led strip
        config = settings.get()
        self.strip = PixelStrip(
            num=config.LED_COUNT,
            pin=config.LED_GPIO_PIN,
            strip_type=config.LED_STRIP_TYPE
        )
        self.strip.begin()

        # Set initial color and brightness
        self._color = BLACK
        self._brightness = 0
        self.update()

        # Init sunrise
        self._sunrise = False
        self._update_sunrise_event = None

        # Register cleanup on exit
        atexit.register(self.cleanup)

    @property
    def color(self):
        return self._color

    def update(self):
        """Set all leds to the current color"""
        for led_index in range(self.strip.numPixels()):
            self.strip.setPixelColorRGB(led_index, *self.color)
        self.strip.setBrightness(self._brightness)
        self.strip.show()

    def cleanup(self):
        """Cleanup on exit"""
        self.set_black()

    def set_color(self, color: Color, brightness: int = 100):
        """Set all leds to a specific color and brightness (0 - 255)"""
        self._color = color
        self._brightness = brightness
        self.update()

    def set_rgb(self, red: int, green: int, blue: int, brightness: int = 100):
        """Set all leds to a RGB value and brightness (0 - 255)"""
        self._color = Color(red, green, blue)
        self._brightness = brightness
        self.update()

    def set_black(self):
        """Turn off all leds"""
        self._color = BLACK
        self._brightness = 0
        self.update()

    def start_sunrise_simulation(self):
        """Start simulating a sunrise"""
        # Check if already in sunrise
        if self._sunrise:
            return

        # Set initial state
        self.set_color(RED, 1)

        # Set timer to update sunrise
        config = settings.get()
        loop = asyncio.get_running_loop()
        self._update_sunrise_event = loop.call_later(
            callback=self.update_sunrise_simulation,
            delay=config.LIGHT_INCREASE_DURATION.seconds / 100,
        )

    def update_sunrise_simulation(self):
        """Update the sunrise simulation"""
        # Derive color from brightness
        brightness = self._brightness + 1
        if brightness > 90:
            color = WARM_WHITE
        elif brightness > 60:
            color = YELLOW
        elif brightness > 30:
            color = ORANGE
        else:
            color = RED

        # Update leds
        self.set_color(color, brightness)

        # Cancel updating sunrise at brightness 100
        if brightness >= 100 and self._update_sunrise_event is not None:
            self._update_sunrise_event.cancel()
            self._update_sunrise_event = None

        # Reschedule
        config = settings.get()
        loop = asyncio.get_running_loop()
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
