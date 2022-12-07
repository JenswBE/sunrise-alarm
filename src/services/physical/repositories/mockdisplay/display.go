package mockdisplay

import (
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
	"github.com/rs/zerolog/log"
)

var _ repositories.Display = &MockDisplay{}

type MockDisplay struct{}

func NewMockDisplay() *MockDisplay {
	return &MockDisplay{}
}

func (d *MockDisplay) SetBrightness(brightness byte) error {
	log.Debug().Uint8("new_value", brightness).Msg("MockDisplay.SetBrightness: Setting new brightness")
	return nil
}

func (d *MockDisplay) Close() error {
	log.Debug().Msg("MockDisplay.Close: Closing device")
	return nil
}
