"""This module contains a helper for a temperature and humidity sensor"""

import asyncio
import logging
from collections import namedtuple
from datetime import timedelta
from typing import Callable

from physical.helpers import settings

if not settings.get().MOCK:
    from adafruit_blinka.microcontroller.bcm283x.pin import Pin
    import adafruit_dht

POLL_TIMEOUT = timedelta(seconds=5)


THReading = namedtuple('THReading', ['temperature', 'humidity'])


class TempHumid:
    """Helper class to work with a DHT22 temperature and humidity sensor"""

    def __init__(self, gpio_pin: int):
        if not settings.get().MOCK:
            # Init device
            self._sensor = adafruit_dht.DHT22(Pin(gpio_pin), use_pulseio=False)

            # Init reading timeout
            self._loop = asyncio.get_event_loop()
            self._reading_event = self._loop.call_later(
                callback=self._read_sensor,
                delay=POLL_TIMEOUT.seconds,
            )

    def set_new_reading_callback(self, callback: Callable[[THReading], None]):
        """Sets the callback for a new reading"""
        self._new_reading_callback = callback

    def _read_sensor(self):
        try:
            reading = THReading(self._sensor.temperature,
                                self._sensor.humidity)
            self.self._new_reading_callback(reading)
        except RuntimeError as error:
            pass
        except Exception as error:
            logging.error("Reading DHT22 failed with error: %s", error)

        # Reset reading event
        self._reading_event = self._loop.call_later(
            callback=self._read_sensor,
            delay=POLL_TIMEOUT.seconds,
        )
