package gpiobutton

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories"
	"github.com/stianeikeland/go-rpio/v4"
)

var _ repositories.Button = &GPIOButton{}

type GPIOButton struct {
	pin         rpio.Pin
	activeState rpio.State
}

func NewGPIOButton(pinNumber int, highIsActive bool) (*GPIOButton, error) {
	if err := rpio.Open(); err != nil {
		return nil, err
	}

	pin := rpio.Pin(pinNumber)
	pin.Input()

	button := &GPIOButton{
		pin:         pin,
		activeState: rpio.High,
	}
	if !highIsActive {
		button.activeState = rpio.Low
	}
	return button, nil
}

func (b *GPIOButton) IsPressed() bool {
	return b.pin.Read() == b.activeState
}
