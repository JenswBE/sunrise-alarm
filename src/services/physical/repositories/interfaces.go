package repositories

import (
	"github.com/JenswBE/sunrise-alarm/src/services/physical/entities"
)

type Button interface {
	IsPressed() bool
}

type Buzzer interface {
	On()
	Off()
}

type Leds interface {
	GetColorAndBrightness() (entities.PresetColor, byte)
	SetColorAndBrightness(color entities.PresetColor, brightness byte)
	Close() error
}

type LightSensor interface {
	GetVisibleLight() (uint32, error)
	Close() error
}
