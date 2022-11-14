package usecases

import (
	"sync"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/services/audio"
	"github.com/JenswBE/sunrise-alarm/src/services/physical"
	"github.com/JenswBE/sunrise-alarm/src/utils/trigger"
	"github.com/rs/zerolog/log"
)

type Ringer struct {
	ringingMutex        sync.Mutex
	ringing             bool
	ringingChan         chan bool
	ringingSinceMinutes uint
	lightDuration       time.Duration
	soundDuration       time.Duration
	abortAlarm          chan<- struct{}
	audioService        audio.Service
	physicalService     physical.Service
}

const tickDuration = time.Minute

func NewRinger(physicalService physical.Service, audioService audio.Service, lightDuration, soundDuration time.Duration, abortAlarm chan<- struct{}) *Ringer {
	// Init
	ringer := &Ringer{
		ringingChan:     make(chan bool, 1),
		lightDuration:   lightDuration,
		soundDuration:   soundDuration,
		abortAlarm:      abortAlarm,
		audioService:    audioService,
		physicalService: physicalService,
	}

	go func(ringingChan chan bool, physicalService physical.Service, audioService audio.Service) {
		trigger := trigger.NewDelayedTrigger()
		for {
			select {
			case shouldRing := <-ringingChan:
				if shouldRing {
					ringer.ringingSinceMinutes = 0
					physicalService.StartSunriseSimulation()
					trigger.Schedule(tickDuration)
				} else {
					physicalService.StopSunriseSimulation()
					physicalService.StopBuzzer()
					if err := audioService.StopMusic(); err != nil {
						log.Error().Err(err).Msg("Ringer: Failed to stop music")
					}
				}
			case <-trigger.C:
				if !ringer.ringing {
					log.Debug().Msg("Ringer.loop: Received delayed trigger, but ringer is not ringing. Ignoring delayed trigger.")
					continue
				}
				ringer.ringingSinceMinutes++
				ringer.handleNextStep()
				trigger.Schedule(tickDuration)
			}
		}
	}(ringer.ringingChan, physicalService, audioService)

	return ringer
}

func (r *Ringer) handleNextStep() {
	// Init
	currentDelay := time.Duration(r.ringingSinceMinutes) * time.Minute
	logger := log.With().Stringer("current_delay", currentDelay).Stringer("sound_duration", r.soundDuration).Stringer("light_duration", r.lightDuration).Logger()
	logger.Debug().Msgf("Ringer.handleNextStep: Handle next alarm step at minute %d", r.ringingSinceMinutes)

	// Check for abort
	abortDelay := r.lightDuration + 10*time.Minute
	if currentDelay > abortDelay {
		logger.Warn().Stringer("abort_delay", abortDelay).Msg("Ringer.handleNextStep: Ringer reached abort limit. Requesting manager to abort alarm.")
		r.abortAlarm <- struct{}{}
		return
	}

	// Check for buzzer
	if currentDelay > r.lightDuration {
		logger.Debug().Msg("Ringer.handleNextStep: Starting buzzer")
		r.physicalService.StartBuzzer()
	}

	// Check for music
	soundDelay := r.lightDuration - r.soundDuration
	if soundDelay.Truncate(time.Minute) == currentDelay.Truncate(time.Minute) {
		logger.Debug().Msg("Ringer.handleNextStep: Starting alarm music")
		if err := r.audioService.PlayMusic(); err != nil {
			logger.Error().Msg("Ringer.handleNextStep: Failed to play music")
		}
	} else if currentDelay > soundDelay {
		logger.Debug().Msg("Ringer.handleNextStep: Increasing alarm volume")
		if err := r.audioService.IncreaseVolume(); err != nil {
			logger.Error().Msg("Ringer.handleNextStep: Failed to increase alarm volume")
		}
	} else {
		logger.Debug().Msg("Ringer.handleNextStep: No action required (yet) for music")
	}
}

func (r *Ringer) Start() {
	log.Debug().Msg("Ringer.Start called")
	r.setRinging(true)
}

func (r *Ringer) Stop() {
	log.Debug().Msg("Ringer.Stop called")
	r.setRinging(false)
}

func (r *Ringer) setRinging(value bool) {
	r.ringingMutex.Lock()
	defer r.ringingMutex.Unlock()
	if r.ringing != value {
		r.ringingChan <- value
		r.ringing = value
	}
}
