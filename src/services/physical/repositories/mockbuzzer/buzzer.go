package mockbuzzer

import (
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
	"github.com/rs/zerolog/log"
)

var _ repositories.Buzzer = &MockBuzzer{}

type MockBuzzer struct{}

func NewMockBuzzer() *MockBuzzer {
	return &MockBuzzer{}
}

func (b *MockBuzzer) On() {
	log.Debug().Str("state", "on").Msg("MockBuzzer: Turning buzzer ON")
}

func (b *MockBuzzer) Off() {
	log.Debug().Str("state", "off").Msg("MockBuzzer: Turning buzzer OFF")
}
