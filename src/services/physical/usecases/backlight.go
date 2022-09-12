package usecases

import "github.com/rs/zerolog/log"

func (s *PhysicalService) LockBacklightBrightness() {
	log.Error().Msg("LockBacklightBrightness called, but not implemented")
}

func (s *PhysicalService) UnlockBacklightBrightness() {
	log.Error().Msg("UnlockBacklightBrightness called, but not implemented")
}
