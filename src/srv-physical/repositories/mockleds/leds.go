package mockleds

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories"
	"github.com/rs/zerolog/log"
)

var _ repositories.Leds = &MockLeds{}

type MockLeds struct {
	currentColor      entities.PresetColor
	currentBrightness byte
}

func (l *MockLeds) GetColorAndBrightness() (entities.PresetColor, byte) {
	log.Debug().Interface("color", l.currentColor).Uint8("brightness", l.currentBrightness).Msg("MockLeds: Current color and brightness requested")
	return l.currentColor, l.currentBrightness
}

func (l *MockLeds) SetColorAndBrightness(color entities.PresetColor, brightness byte) {
	log.Debug().Interface("color", color).Uint8("brightness", brightness).Msg("MockLeds: Setting new color and brightness")
	l.currentColor = color
	l.currentBrightness = brightness
}