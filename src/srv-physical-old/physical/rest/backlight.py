"""This module handles all routes for Backlight operations"""

from fastapi import APIRouter, Depends

from physical.devices.common import Devices, dev_from_req


router = APIRouter()


@router.put("/lock", summary="Lock the backlight brightness")
async def lock(devs: Devices = Depends(dev_from_req)):
    devs.display.lock_brightness()
    return {}


@router.delete("/lock", summary="Unlock the backlight brightness")
async def unlock(devs: Devices = Depends(dev_from_req)):
    devs.display.unlock_brightness()
    return {}
