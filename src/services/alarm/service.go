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
	// GetNextRingTime returns the next time an alarm will ring.
	// If there are no future alarms, a zero time.Time will be returned.
	GetNextRingTime() time.Time
	CreateAlarm(alarm entities.Alarm) (entities.Alarm, error)
	UpdateAlarm(alarm entities.Alarm) error
	DeleteAlarm(id uuid.UUID) error
	Close() error
}

func NewAlarmService(physicalService physical.Service, audioService audio.Service, pubSub pubsub.PubSub) (Service, error) {
	return usecases.NewAlarmService(physicalService, audioService, pubSub)
}
