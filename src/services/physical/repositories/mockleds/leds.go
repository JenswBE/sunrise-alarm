package mockleds

import (
	"github.com/JenswBE/sunrise-alarm/src/services/physical/entities"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
	"github.com/rs/zerolog/log"
)

var _ repositories.Leds = &MockLeds{}

type MockLeds struct {
	currentColor      entities.PresetColor
	currentBrightness byte
}

func NewMockLeds() *MockLeds {
	return &MockLeds{}
}

func (l *MockLeds) GetColorAndBrightness() (entities.PresetColor, byte) {
	log.Debug().Stringer("color", l.currentColor).Uint8("brightness", l.currentBrightness).Msg("MockLeds.GetColorAndBrightness: Current color and brightness requested")
	return l.currentColor, l.currentBrightness
}

func (l *MockLeds) Close() error {
	log.Debug().Msg("MockLeds: Closing device")
	return nil
}

func (l *MockLeds) SetColorAndBrightness(color entities.PresetColor, brightness byte) {
	log.Debug().Stringer("color", color).Uint8("brightness", brightness).Msg("MockLeds.SetColorAndBrightness: Setting new color and brightness")
	l.currentColor = color
	l.currentBrightness = brightness
}
