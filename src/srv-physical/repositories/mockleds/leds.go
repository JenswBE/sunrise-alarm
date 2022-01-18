package mockleds

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories"
)

var _ repositories.Leds = &MockLeds{}

type MockLeds struct {
	currentColor      entities.Color
	currentBrightness byte
}

func (l *MockLeds) GetColorAndBrightness() (entities.Color, byte) {
	return l.currentColor, l.currentBrightness
}

func (l *MockLeds) SetColorAndBrightness(color entities.Color, brightness byte) {
	l.currentColor = color
	l.currentBrightness = brightness
}
