package gpiobutton

import (
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
	"github.com/stianeikeland/go-rpio/v4"
)

var _ repositories.Button = &GPIOButton{}

type GPIOButton struct {
	pin         rpio.Pin
	activeState rpio.State
}

func NewGPIOButton(pinNumber int, highIsActive bool) *GPIOButton {
	pin := rpio.Pin(pinNumber)
	pin.Input()

	button := &GPIOButton{pin: pin}
	if highIsActive {
		button.activeState = rpio.High
		button.pin.PullDown()
	} else {
		button.activeState = rpio.Low
		button.pin.PullUp()
	}
	return button
}

func (b *GPIOButton) IsPressed() bool {
	return b.pin.Read() == b.activeState
}
