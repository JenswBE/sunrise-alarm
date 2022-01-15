"""This module handles all routes for LEDs operations"""

from fastapi import APIRouter, Depends
from pydantic import BaseModel, conint

from physical.devices import leds
from physical.devices.common import Devices, dev_from_req


router = APIRouter()


class Leds(BaseModel):
    color: leds.PresetColor
    brightness: conint(ge=0, le=255) = 100


@router.get("/", summary="Get current color and brightness of LEDs", response_model=Leds)
async def get_leds(devs: Devices = Depends(dev_from_req)):
    return Leds(
        color=devs.leds.color,
        brightness=devs.leds.brightness,
    )


@router.put("/", summary="Set LEDs to specified color and brightness")
async def set_leds(req: Leds, devs: Devices = Depends(dev_from_req)):
    devs.leds.set_color(req.color, req.brightness)
    return {}


@router.delete("/", summary="Turn off all LEDs")
async def clear_leds(devs: Devices = Depends(dev_from_req)):
    devs.leds.set_black()
    return {}


@router.put("/sunrise", summary="Start sunrise simulation")
async def start_sunrise(devs: Devices = Depends(dev_from_req)):
    devs.leds.start_sunrise_simulation()
    return {}


@router.delete("/sunrise", summary="Stop sunrise simulation")
async def stop_sunrise(devs: Devices = Depends(dev_from_req)):
    devs.leds.stop_sunrise_simulation()
    return {}
