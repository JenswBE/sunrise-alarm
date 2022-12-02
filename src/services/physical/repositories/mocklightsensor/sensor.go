package mocklightsensor

import (
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
	"github.com/rs/zerolog/log"
)

var _ repositories.LightSensor = &MockLightSensor{}

type MockLightSensor struct {
	visibleLight uint32
}

func NewMockLightSensor(visibleLight uint32) *MockLightSensor {
	return &MockLightSensor{visibleLight: visibleLight}
}

func (s *MockLightSensor) GetVisibleLight() (uint32, error) {
	log.Debug().Uint32("value", s.visibleLight).Msg("MockLightSensor.GetVisibleLight: Getting visible light")
	return s.visibleLight, nil
}

func (s *MockLightSensor) SetVisibleLight(visibleLight uint32) {
	log.Debug().Uint32("new_value", visibleLight).Msg("MockLightSensor.SetVisibleLight: Setting visible light")
	s.visibleLight = visibleLight
}

func (s *MockLightSensor) Close() error {
	log.Debug().Msg("MockLightSensor.Close: Closing device")
	return nil
}
