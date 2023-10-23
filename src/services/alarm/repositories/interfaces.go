package repositories

import (
	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/google/uuid"
)

type DB interface {
	List() ([]entities.Alarm, error)
	Get(id uuid.UUID) (entities.Alarm, error)
	Create(alarm entities.Alarm) (entities.Alarm, error)
	Update(alarm entities.Alarm) error
	Delete(id uuid.UUID) error
	Close() error
}
