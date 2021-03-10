"""This module handles publishing and receiving of MQTT messages."""

import json
import logging
import random

from gmqtt import Client
from gmqtt.mqtt.constants import MQTTv311
from starlette.requests import Request

from physical.helpers import settings
from physical.devices.temp_humid import THReading

mqtt = False


class MQTT:
    """Helper class to abstract MQTT client implementation"""
    _topic_prefix = "sunrise_alarm/"

    def __init__(self, client):
        self.client = client

    async def stop(self):
        await self.client.disconnect()

    def _publish(self, topic: str, payload: str):
        full_topic = self._topic_prefix + topic
        logging.info("Publish event to MQTT topic: %s", full_topic)
        self.client.publish(full_topic, payload, qos=1)

    def publish_button_pressed(self):
        self._publish("button_pressed", "")

    def publish_button_long_pressed(self):
        self._publish("button_long_pressed", "")

    def publish_temp_humid_updated(self, reading: THReading):
        self._publish("temp_humid_updated", reading.json())


async def get() -> MQTT:
    """Returns current or creates new MQTT helper"""
    global mqtt
    if mqtt is not False:
        return mqtt

    # Build new client from settings
    config = settings.get()
    client_id = "{}-{:08x}".format(config.MQTT_CLIENT_ID,
                                   random.randrange(2**32))
    client = Client(client_id)
    await client.connect(config.MQTT_BROKER_HOST, config.MQTT_BROKER_PORT, keepalive=60, version=MQTTv311)
    mqtt = MQTT(client)
    return mqtt


def mqtt_from_req(request: Request) -> MQTT:
    """Returns instance of MQTT from request"""
    return request.app.state.mqtt
