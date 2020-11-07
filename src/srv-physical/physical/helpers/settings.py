"""This module handles all ENV and file based settings."""

import logging
from functools import lru_cache

from pydantic import BaseSettings


class Settings(BaseSettings):
    """Handles ENV and file based settings"""
    # General
    DEBUG: bool = False
    MOCK: bool = False

    # MQTT
    MQTT_HOST: str = "localhost"
    MQTT_PORT: int = 1883
    MQTT_CLIENT_ID: str = "srv-physical"

    # BUTTON
    BUTTON_GPIO_PIN: int = 23


@lru_cache(maxsize=None)
def get():
    """Returns cached settings object"""
    settings = Settings()

    if settings.DEBUG:
        logging.getLogger().setLevel(logging.DEBUG)

    return settings
