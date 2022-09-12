package alarm

import (
	"time"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/JenswBE/sunrise-alarm/src/services/alarm/usecases"
	"github.com/JenswBE/sunrise-alarm/src/services/audio"
	"github.com/JenswBE/sunrise-alarm/src/services/physical"
	"github.com/JenswBE/sunrise-alarm/src/utils/pubsub"
	"github.com/google/uuid"
)

type Service interface {
	ListAlarms() ([]entities.Alarm, error)
	GetAlarm(id uuid.UUID) (entities.Alarm, error)
	// GetNextAlarmToRing returns the next alarm that will ring.
	// If there are no future alarms, nil will be returned.
	GetNextAlarmToRing() *entities.NextAlarm
	CreateAlarm(alarm entities.Alarm) (entities.Alarm, error)
	UpdateAlarm(alarm entities.Alarm) error
	DeleteAlarm(id uuid.UUID) error
	Close() error
}

func NewAlarmService(physicalService physical.Service, audioService audio.Service, pubSub pubsub.PubSub, alarmLightDuration time.Duration, alarmSoundDuration time.Duration) Service {
	return usecases.NewAlarmService(physicalService, audioService, pubSub, alarmLightDuration, alarmSoundDuration)
}
