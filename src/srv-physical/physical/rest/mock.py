"""This module handles all routes for mock operations"""

from fastapi import APIRouter, Depends

from physical.helpers.mqtt import MQTT, mqtt_from_req


router = APIRouter()


@router.post("/button/pressed", summary="Mock press on the button")
async def mock_button_pressed(mqtt: MQTT = Depends(mqtt_from_req)):
    """Mock press on the button"""
    mqtt.publish_button_pressed()
    return {}


@router.post("/button/long_pressed", summary="Mock long press on the button")
async def mock_button_long_pressed(mqtt: MQTT = Depends(mqtt_from_req)):
    """Mock long press on the button"""
    mqtt.publish_button_long_pressed()
    return {}
