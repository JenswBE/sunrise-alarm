package usecases

import (
	"github.com/rs/zerolog/log"
)

func (s *PhysicalService) StartBuzzer() {
	log.Debug().Msg("Starting buzzer sequencer")
	s.seq.Start()
}

func (s *PhysicalService) StopBuzzer() {
	log.Debug().Msg("Stopping buzzer sequencer")
	s.seq.Stop()
}
