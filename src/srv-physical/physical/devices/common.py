from dataclasses import dataclass

from starlette.requests import Request

from physical.devices.button import Button
from physical.devices.buzzer import Buzzer
from physical.devices.display import Display
from physical.devices.leds import Leds


@dataclass
class Devices:
    button: Button
    buzzer: Buzzer
    display: Display
    leds: Leds


def dev_from_req(request: Request) -> Devices:
    """Returns devices from request"""
    return request.app.state.devices
