package repositories

import (
	"context"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
)

type Button interface {
	IsPressed() bool
}

type Leds interface {
	GetColorAndBrightness() (entities.PresetColor, byte)
	SetColorAndBrightness(color entities.PresetColor, brightness byte)
}

type MQTTClient interface {
	Publish(ctx context.Context, topic, payload string) error
}
