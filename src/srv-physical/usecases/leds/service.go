package leds

import (
	"sync"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories"
	"github.com/rs/zerolog/log"
)

var _ Usecase = &Service{}

type Service struct {
	leds            repositories.Leds
	sunriseStop     chan bool
	sunriseLock     sync.Mutex // Ensures we don't mess with sunriseStop between nil check and sunriseStop creation/deletion
	sunriseDuration time.Duration
}

func NewService(leds repositories.Leds, sunriseDuration time.Duration) *Service {
	return &Service{
		leds:            leds,
		sunriseDuration: sunriseDuration,
	}
}

func (s *Service) GetColorAndBrightness() (entities.PresetColor, byte) {
	log.Debug().Msg("Leds Service: Getting current leds color and brightness")
	return s.leds.GetColorAndBrightness()
}

func (s *Service) SetColorAndBrightness(color entities.PresetColor, brightness byte) {
	log.Debug().Interface("color", color).Uint8("brightness", brightness).Msg("Leds Service: Setting new color and brightness")
	s.leds.SetColorAndBrightness(color, brightness)
}

func (s *Service) Clear() {
	log.Debug().Msg("Leds Service: Clearing leds")
	s.leds.SetColorAndBrightness(entities.PresetColorBlack, 0)
}

// Start simulating a sunrise
func (s *Service) StartSunrise() {
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
func (s *Service) StopSunrise() {
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
	s.Clear()
}

func (s *Service) runSunrise(stop chan bool) {
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
			s.SetColorAndBrightness(color, brightness)
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
