package leds

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
)

type Usecase interface {
	GetColorAndBrightness() (entities.PresetColor, byte)
	SetColorAndBrightness(color entities.PresetColor, brightness byte)
	Clear()
}
