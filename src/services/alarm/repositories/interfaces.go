package repositories

import (
	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/google/uuid"
)

type DB interface {
	List() ([]entities.Alarm, error)
	Get(id uuid.UUID) (entities.Alarm, error)
	Create(entities.Alarm) (entities.Alarm, error)
	Update(entities.Alarm) error
	Delete(id uuid.UUID) error
	Close() error
}
