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
        self._button.when_pressed = self.handle_press

        # Setup other variables
        self._first_press_timestamp = None
        self._is_long_press = False

    def set_short_press_callback(self, callback):
        """Sets the callback for a short press"""
        self._short_press_callback = callback

    def set_long_press_callback(self, callback):
        """Sets the callback for a long press"""
        self._long_press_callback = callback

    def handle_press(self):
        """Handle the button press callback"""
        # Ignore if we are in a long press
        if self._is_long_press:
            return

        # Debounce button
        if self._first_press_timestamp is not None and datetime.now() - self._first_press_timestamp < DEBOUNCE_DURATION:
            return

        # Save timestamp of first press
        if not self._first_press_timestamp:
            self._first_press_timestamp = datetime.now()

        # Schedule a check for the press length
        loop = asyncio.get_running_loop()
        loop.call_later(0.1, self.check_press_length)

    def check_press_length(self):
        """Check if it's a short or long press"""
        # Check if button is still pressed
        if self._button.is_pressed:
            # Schedule a new check
            loop = asyncio.get_running_loop()
            loop.call_later(0.1, self.check_press_length)

            # Handle edge case due to multiple clicks
            if self._first_press_timestamp is None:
                return

            # Check if we reached a long press
            diff = datetime.now() - self._first_press_timestamp
            if not self._is_long_press and diff > LONG_PRESS_DURATION:
                # Long press reached
                self._is_long_press = True
                if self._long_press_callback is not None:
                    self._long_press_callback()

            # Long press handled
            return

        # Handle short press
        if not self._is_long_press:
            if self._short_press_callback is not None:
                self._short_press_callback()

        # Clean state on button released
        self._first_press_timestamp = None
        self._is_long_press = False
