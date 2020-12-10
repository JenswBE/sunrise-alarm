"""This module contains helpers for physical buttons"""

import asyncio
from datetime import datetime, timedelta

from physical.helpers import settings

if not settings.get().MOCK:
    from gpiozero import Button as GPIOButton

DEBOUNCE_DURATION = timedelta(seconds=0.2)
LONG_PRESS_DURATION = timedelta(seconds=1)


class Button:
    """
    Helper class to work with a physical button

    Debounce and on_held of gpiozero is not working as expected.
    Therefore, I implemented this custom helper class.
    """

    def __init__(self, gpio_pin: int):
        # Setup button
        self._button = GPIOButton(pin=gpio_pin, pull_up=False)
        self._button.when_held = self.handle_held
        self._button.when_released = self.handle_release
        self._was_held = False

    def set_short_press_callback(self, callback):
        """Sets the callback for a short press"""
        self._short_press_callback = callback

    def set_long_press_callback(self, callback):
        """Sets the callback for a long press"""
        self._long_press_callback = callback

    def handle_held(self):
        self._was_held = True
        if self._long_press_callback is not None:
            self._long_press_callback()

    def handle_release(self):
        if not self._was_held and self._short_press_callback is not None:
            self._short_press_callback()
        self._was_held = False
