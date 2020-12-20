"""This module handles all ENV and file based settings."""

import logging
from datetime import timedelta
from functools import lru_cache

from pydantic import BaseSettings
from rpi_ws281x import ws


class Settings(BaseSettings):
    """Handles ENV and file based settings"""
    # General
    MOCK: bool = False

    # MQTT
    MQTT_BROKER_HOST: str = "localhost"
    MQTT_BROKER_PORT: int = 1883
    MQTT_CLIENT_ID: str = "srv-physical"

    # BUTTON
    BUTTON_GPIO_PIN: int = 23

    # LEDS
    # See https://github.com/rpi-ws281x/rpi-ws281x-python
    LED_STRIP_TYPE: int = ws.WS2811_STRIP_GRB
    LED_COUNT: int = 33
    LED_GPIO_PIN: int = 21
    LIGHT_INCREASE_DURATION: timedelta = timedelta(minutes=5)


@lru_cache(maxsize=None)
def get():
    """Returns cached settings object"""
    return Settings()
