package leds

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
)

type Usecase interface {
	GetColorAndBrightness() (entities.Color, byte)
	SetColorAndBrightness(color entities.Color, brightness byte)
	Clear()
}
