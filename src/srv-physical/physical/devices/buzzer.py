"""This module contains helpers for physical buttons"""

import asyncio
import logging
from datetime import datetime, timedelta

from physical.helpers import settings

if not settings.get().MOCK:
    from gpiozero import Buzzer as GPIOBuzzer

BEEP = [
    timedelta(seconds=.1),
    timedelta(seconds=.1),
    timedelta(seconds=.1),
    timedelta(seconds=1),
]


class Buzzer:
    """Helper class to work with a buzzer"""

    def __init__(self, gpio_pin: int):
        self._buzzer = GPIOBuzzer(pin=gpio_pin)
        self._loop = asyncio.get_running_loop()
        self._enabled = False

    def start(self, callback):
        """Starts the buzzer"""
        self._enabled = True
        self._beep_step = 0
        self._beep_step_event = self._loop.call_later(
            callback=self.handle_buzzer_step,
            delay=BEEP[self._beep_step].seconds,
        )

    def stop(self, callback):
        """Stops the buzzer"""
        self._enabled = False

    def handle_buzzer_step(self):
        """Handle a step in the buzzer sequence"""
        if self._enabled:
            # Update buzzer step
            self._beep_step += 1
            if self._beep_step >= len(BEEP):
                self._beep_step = 0

            # Update buzzer
            if self._beep_step % 2:
                self._buzzer.on()
            else:
                self._buzzer.off()

            # Schedule new buzzer beep
            self._beep_step_event = self._loop.call_later(
                callback=self.handle_buzzer_step,
                delay=BEEP[self._beep_step].seconds,
            )
