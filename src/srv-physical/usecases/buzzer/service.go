package buzzer

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/utils/buzzersequencer"
	"github.com/rs/zerolog/log"
)

var _ Usecase = &Service{}

type Service struct {
	seq *buzzersequencer.BuzzerSequencer
}

func NewService(buzzer repositories.Buzzer) *Service {
	seq := buzzersequencer.NewBuzzerSequencer(buzzer)
	return &Service{seq: seq}
}

func (s *Service) Start() {
	log.Debug().Msg("Buzzer Service: Starting buzzer sequencer")
	s.seq.Start()
}

func (s *Service) Stop() {
	log.Debug().Msg("Buzzer Service: Stopping buzzer sequencer")
	s.seq.Stop()
}
