package repositories

import (
	"context"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
)

type Button interface {
	IsPressed() bool
}

type Leds interface {
	GetColorAndBrightness() (entities.Color, byte)
	SetColorAndBrightness(color entities.Color, brightness byte)
}

type MQTTClient interface {
	Publish(ctx context.Context, topic, payload string) error
}
