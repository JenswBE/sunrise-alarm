package physical

import (
	"github.com/JenswBE/sunrise-alarm/src/cmd/config"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/entities"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/usecases"
	"github.com/JenswBE/sunrise-alarm/src/utils/pubsub"
)

type Service interface {
	LockBacklightBrightness()
	UnlockBacklightBrightness()

	StartBuzzer()
	StopBuzzer()

	GetLEDState() entities.LEDState
	SetLEDState(state entities.LEDState)
	ResetLEDState()
	StartSunriseSimulation()
	StopSunriseSimulation()

	Close()
}

func NewPhysicalService(config config.PhysicalConfig, pubSub pubsub.PubSub) Service {
	return usecases.NewPhysicalService(config, pubSub)
}
