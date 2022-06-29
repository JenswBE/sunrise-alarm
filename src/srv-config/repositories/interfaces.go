package repositories

import (
	"context"

	"github.com/JenswBE/sunrise-alarm/src/srv-config/entities"
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

type MQTTClient interface {
	Publish(ctx context.Context, topic, payload string) error
}
