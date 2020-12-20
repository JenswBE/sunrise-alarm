"""This module handles publishing and receiving of MQTT messages."""

import asyncio
import functools
import logging
import random

from gmqtt import Client
from gmqtt.mqtt.constants import MQTTv311
from starlette.requests import Request

from physical.helpers import settings

mqtt = False


class MQTT:
    """Helper class to abstract MQTT client implementation"""
    _topic_prefix = "sunrise_alarm/"

    def __init__(self, client):
        self.client = client

    async def stop(self):
        await self.client.disconnect()

    def _publish(self, topic):
        full_topic = self._topic_prefix + topic
        logging.info("Publish event to MQTT topic: %s", full_topic)
        self.client.publish(full_topic, qos=1)

    def publish_button_pressed(self):
        self._publish("button_pressed")

    def publish_button_long_pressed(self):
        self._publish("button_long_pressed")


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
