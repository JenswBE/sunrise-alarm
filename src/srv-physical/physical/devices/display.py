"""This module contains helpers for the display backlight"""

import asyncio
import atexit
import logging
from datetime import timedelta

from physical.helpers import settings

if not settings.get().MOCK:
    import board
    import busio
    from adafruit_tsl2591 import TSL2591
    from rpi_backlight import Backlight

MIN_LIGHT = 66000
MAX_LIGHT = 1700000
MIN_BRIGHTNESS = 3
MAX_BRIGHTNESS = 70
SLEEP_TIMEOUT = timedelta(seconds=10)


class Display:
    """Helper class to work with the display backlight"""

    def __init__(self):
        # Init light sensor
        i2c = busio.I2C(board.SCL, board.SDA)
        self._sensor = TSL2591(i2c)

        # Init backlight
        self._backlight = Backlight()
        self._backlight.fade_duration = 0.25

        # Init sleep timeout
        self._backlight.power = True
        loop = asyncio.get_running_loop()
        self._sleep_event = loop.call_later(
            callback=self.sleep,
            delay=SLEEP_TIMEOUT.seconds,
        )
        atexit.register(self.enable_keep_awake)

        # Init remaining variables
        self._brightness_locked = False
        self._sleep_callback = None

        # Start updating brightness
        loop.call_later(1, self.update_brightness)

    def lock_brightness(self):
        """Lock the backlight brightness"""
        self._brightness_locked = True

    def unlock_brightness(self):
        """Unlock the backlight brightness"""
        # Delay unlock as sensor caches readings
        loop = asyncio.get_running_loop()
        loop.call_later(1, self._unlock_brightness)

    def _unlock_brightness(self):
        """Unlock the backlight brightness (really this time)"""
        self._brightness_locked = False

    def update_brightness(self):
        """Update the current backlight brightness"""
        # Check if brightness is locked
        if self._brightness_locked:
            return

        # Update brightness
        current_light = self._sensor.visible
        new_brightness = self.calculate_brightness(current_light)
        self._backlight.brightness = new_brightness

        # Reschedule call
        loop = asyncio.get_running_loop()
        loop.call_later(1, self.update_brightness)

    def calculate_brightness(self, current_light):
        """Calculates the current brightness for the display"""
        # Check if minimum light is reached
        if current_light <= MIN_LIGHT:
            return MIN_BRIGHTNESS

        # Check if maximum light is reached
        if current_light >= MAX_LIGHT:
            return MAX_BRIGHTNESS

        # Calculate new brightness
        # See https://stackoverflow.com/a/929107
        light_range = MAX_LIGHT - MIN_LIGHT
        brightness_range = MAX_BRIGHTNESS - MIN_BRIGHTNESS
        new_brightness = (((current_light - MIN_LIGHT) *
                           brightness_range) / light_range) + MIN_BRIGHTNESS

        # Strip decimals and return result
        return int(new_brightness)

    def sleep(self):
        """Put screen into sleep"""
        # Disable screen
        logging.info("Display: Put display into sleep")
        self._backlight.power = False

        # Call callback if set
        if self._sleep_callback is not None:
            self._sleep_callback()

    def set_sleep_callback(self, callback):
        """Set a callback to be called on sleep"""
        self._sleep_callback = callback

    def wake(self):
        """Wakes the screen"""
        # Check if screen is in keep_awake
        if self._backlight.power and self._sleep_event is None:
            return

        # (Re)set timeout
        if self._sleep_event is not None:
            self._sleep_event.cancel()
        loop = asyncio.get_running_loop()
        self._sleep_event = loop.call_later(
            callback=self.sleep,
            delay=SLEEP_TIMEOUT.seconds,
        )

        # Enable backlight
        logging.info("Display: Awake display")
        self._backlight.power = True

    def enable_keep_awake(self):
        """Keep screen awake until next sleep"""
        # Cancel timeout if any
        if self._sleep_event is not None:
            self._sleep_event.cancel()
            self._sleep_event = None

        # Enable backlight
        logging.info("Display: Keep display awake")
        self._backlight.power = True

    def disable_keep_awake(self):
        """Allow screen to go to sleep again"""
        # Ignore call if screen is already sleeping
        if not self._backlight.power:
            return

        # Log action
        logging.info("Display: Allow display to go to sleep")

        # (Re)set timeout
        if self._sleep_event is not None:
            self._sleep_event.cancel()
        loop = asyncio.get_running_loop()
        self._sleep_event = loop.call_later(
            callback=self.sleep,
            delay=SLEEP_TIMEOUT.seconds,
        )
