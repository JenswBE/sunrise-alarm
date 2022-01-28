package leds

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories"
	"github.com/rs/zerolog/log"
)

var _ Usecase = &Service{}

type Service struct {
	leds repositories.Leds
}

func NewService(leds repositories.Leds) *Service {
	return &Service{leds: leds}
}

func (s *Service) GetColorAndBrightness() (entities.PresetColor, byte) {
	log.Debug().Msg("Leds Service: Getting current leds color and brightness")
	return s.leds.GetColorAndBrightness()
}

func (s *Service) SetColorAndBrightness(color entities.PresetColor, brightness byte) {
	log.Debug().Interface("color", color).Uint8("brightness", brightness).Msg("Leds Service: Setting new color and brightness")
	s.leds.SetColorAndBrightness(color, brightness)
}

func (s *Service) Clear() {
	log.Debug().Msg("Leds Service: Clearing leds")
	s.leds.SetColorAndBrightness(entities.PresetColorBlack, 0)
}
