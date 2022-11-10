package usecases

import (
	"time"

	"github.com/JenswBE/sunrise-alarm/src/services/physical/entities"
	"github.com/rs/zerolog/log"
)

func (s *PhysicalService) GetLEDState() entities.LEDState {
	log.Debug().Msg("Leds Service: Getting current leds color and brightness")
	color, brightness := s.leds.GetColorAndBrightness()
	return entities.LEDState{
		Color:      color,
		Brightness: brightness,
	}
}

func (s *PhysicalService) SetLEDState(state entities.LEDState) {
	log.Debug().Stringer("color", state.Color).Uint8("brightness", state.Brightness).Msg("Leds Service: Setting new color and brightness")
	s.leds.SetColorAndBrightness(state.Color, state.Brightness)
}

func (s *PhysicalService) ResetLEDState() {
	log.Debug().Msg("Leds Service: Resetting leds")
	s.leds.SetColorAndBrightness(entities.PresetColorBlack, 0)
}

// Start simulating a sunrise
func (s *PhysicalService) StartSunriseSimulation() {
	log.Debug().Msg("Leds Service: Starting sunrise")
	s.sunriseLock.Lock()
	defer s.sunriseLock.Unlock()
	if s.sunriseStop != nil {
		log.Info().Msg("Leds Service: We are already in a sunrise. Ignoring request to start another one.")
		return
	}

	// Start sunrise
	s.sunriseStop = make(chan bool)
	go s.runSunrise(s.sunriseStop)
}

// Stop simulating a sunrise
func (s *PhysicalService) StopSunriseSimulation() {
	log.Debug().Msg("Leds Service: Stopping sunrise")
	s.sunriseLock.Lock()
	defer s.sunriseLock.Unlock()
	if s.sunriseStop == nil {
		log.Info().Msg("Leds Service: We are not in a sunrise. Ignoring request to stop the sunrise.")
		return
	}

	// Stop sunrise
	close(s.sunriseStop)
	s.sunriseStop = nil
	s.ResetLEDState()
}

func (s *PhysicalService) runSunrise(stop chan bool) {
	var brightness byte = 4 // First brightness will be 5
	for {
		select {
		case <-stop:
			log.Debug().Msg("Leds Service: runSunrise received stop signal")
			return
		default:
			brightness++
			var color entities.PresetColor
			switch {
			case brightness > 90:
				color = entities.PresetColorWarmWhite
			case brightness > 60:
				color = entities.PresetColorYellow
			case brightness > 30:
				color = entities.PresetColorOrange
			default:
				color = entities.PresetColorRed
			}
			log.Debug().Stringer("color", color).Uint8("brightness", brightness).Msg("Leds Service: runSunrise is updating leds")
			s.leds.SetColorAndBrightness(color, brightness)
			if brightness >= 100 {
				log.Debug().Msg("Leds Service: runSunrise reached final state of sunrise. Sustaining color/brightness and waiting for stop signal.")
				<-stop
				log.Debug().Msg("Leds Service: runSunrise received stop signal on final state.")
				return
			}
			time.Sleep(s.sunriseDuration / 95) // 5 to 100 has 95 steps
		}
	}
}
