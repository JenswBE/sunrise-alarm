package autobacklight

import (
	"math"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
	"github.com/rs/zerolog/log"
)

const initialDelay = 10 * time.Second
const updateDelay = time.Second
const minLight = 66000
const maxLight = 1700000
const minBrightness = 0
const maxBrightness = math.MaxUint8

type AutoBacklight struct {
	display     repositories.Display
	lightSensor repositories.LightSensor

	locked bool
	ticker *time.Ticker
}

func NewAutoBacklight(display repositories.Display, lightSensor repositories.LightSensor) *AutoBacklight {
	// Init
	ab := &AutoBacklight{
		display:     display,
		lightSensor: lightSensor,
		ticker:      time.NewTicker(updateDelay),
	}

	// Start updating after initial delay
	time.AfterFunc(initialDelay, func() {
		for range ab.ticker.C {
			if !ab.locked {
				ab.update()
			}
		}
	})
	return ab
}

func (ab *AutoBacklight) update() {
	// Get light
	visibleLight, err := ab.lightSensor.GetVisibleLight()
	if err != nil {
		log.Warn().Err(err).Msg("AutoBacklight.update: Failed to read visible light from sensor. Ignoring error...")
		return
	}

	// Calculate brightness
	var brightness byte
	switch {
	case visibleLight <= minLight:
		brightness = 0
	case visibleLight >= maxLight:
		brightness = math.MaxUint8
	default:
		lightRange := maxLight - minLight
		brightnessRange := maxBrightness - minBrightness
		brightness = byte(float64((visibleLight-minLight)*uint32(brightnessRange))/float64(lightRange) + minBrightness)
	}

	// Set new brightness
	log.Trace().Uint32("visible_light", visibleLight).Uint8("new_brightness", brightness).Msg("AutoBacklight.update: Updating brightness of display backlight...")
	if err := ab.display.SetBrightness(brightness); err != nil {
		log.Warn().Err(err).Msg("AutoBacklight.update: Failed to update display with new brightness. Ignoring error...")
		return
	}
}

func (ab *AutoBacklight) LockBrightness() {
	ab.locked = true
}

func (ab *AutoBacklight) UnlockBrightness() {
	ab.locked = false
}
