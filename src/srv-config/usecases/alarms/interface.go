package alarms

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-config/entities"
	"github.com/google/uuid"
)

type Usecase interface {
	ListAlarms() ([]entities.Alarm, error)
	GetAlarm(id uuid.UUID) (entities.Alarm, error)
	CreateAlarm(alarm entities.Alarm) (entities.Alarm, error)
	UpdateAlarm(alarm entities.Alarm) error
	DeleteAlarm(id uuid.UUID) error
}
