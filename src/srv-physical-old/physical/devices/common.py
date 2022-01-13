from dataclasses import dataclass

from starlette.requests import Request

from physical.devices.button import Button
from physical.devices.buzzer import Buzzer
from physical.devices.display import Display
from physical.devices.leds import Leds
from physical.devices.temp_humid import TempHumid


@dataclass
class Devices:
    button: Button
    buzzer: Buzzer
    display: Display
    leds: Leds
    th_sensor: TempHumid


def dev_from_req(request: Request) -> Devices:
    """Returns devices from request"""
    return request.app.state.devices
