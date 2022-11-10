package gpiobuzzer

import (
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
	"github.com/stianeikeland/go-rpio/v4"
)

var _ repositories.Buzzer = &GPIOBuzzer{}

type GPIOBuzzer struct {
	pin rpio.Pin
}

func NewGPIOBuzzer(pinNumber int) *GPIOBuzzer {
	pin := rpio.Pin(pinNumber)
	pin.Output()
	pin.Low()

	return &GPIOBuzzer{pin: pin}
}

func (b *GPIOBuzzer) On() {
	b.pin.High()
}

func (b *GPIOBuzzer) Off() {
	b.pin.Low()
}
