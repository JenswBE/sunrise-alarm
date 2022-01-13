"""This module handles all routes for Buzzer operations"""

from fastapi import APIRouter, Depends

from physical.devices.common import Devices, dev_from_req


router = APIRouter()


@router.put("/", summary="Start the buzzer")
async def start(devs: Devices = Depends(dev_from_req)):
    devs.buzzer.start()
    return {}


@router.delete("/", summary="Stop the buzzer")
async def stop(devs: Devices = Depends(dev_from_req)):
    devs.buzzer.stop()
    return {}
