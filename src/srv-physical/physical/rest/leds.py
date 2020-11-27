"""This module handles all routes for LEDs operations"""

from fastapi import APIRouter, Depends, status
from fastapi.responses import JSONResponse
from pydantic import BaseModel, conint

from physical.devices import leds
from physical.devices.common import Devices, dev_from_req


router = APIRouter()

RESPONSES = {
    200: {"description": "Success"},
    202: {"description": "Action was mocked"},
}


class SetLedsRequest(BaseModel):
    color: leds.PresetColorEnum
    brightness: conint(ge=0, le=255) = 100


@router.post("/", summary="Set LEDs to specified color and brightness", responses=RESPONSES)
async def set_leds(req: SetLedsRequest, devs: Devices = Depends(dev_from_req)):
    if devs.leds is not None:
        devs.leds.set_color(getattr(leds, req.color), req.brightness)
        return {}
    return JSONResponse(status_code=status.HTTP_202_ACCEPTED, content={})


@router.delete("/", summary="Turn off all LEDs", responses=RESPONSES)
async def clear_leds(devs: Devices = Depends(dev_from_req)):
    if devs.leds is not None:
        devs.leds.set_black()
        return {}
    return JSONResponse(status_code=status.HTTP_202_ACCEPTED, content={})


@router.post("/sunrise", summary="Start sunrise simulation", responses=RESPONSES)
async def start_sunrise(devs: Devices = Depends(dev_from_req)):
    if devs.leds is not None:
        devs.leds.start_sunrise_simulation()
        return {}
    return JSONResponse(status_code=status.HTTP_202_ACCEPTED, content={})


@router.delete("/sunrise", summary="Stop sunrise simulation", responses=RESPONSES)
async def stop_sunrise(devs: Devices = Depends(dev_from_req)):
    if devs.leds is not None:
        devs.leds.stop_sunrise_simulation()
        return {}
    return JSONResponse(status_code=status.HTTP_202_ACCEPTED, content={})
