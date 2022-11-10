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
	managerActions      chan<- ManagerAction
	audioService        audio.Service
	physicalService     physical.Service
}

const tickDuration = time.Minute

func NewRinger(physicalService physical.Service, audioService audio.Service, lightDuration, soundDuration time.Duration, managerActions chan<- ManagerAction) *Ringer {
	// Init
	ringer := &Ringer{
		ringingChan:     make(chan bool, 1),
		lightDuration:   lightDuration,
		soundDuration:   soundDuration,
		managerActions:  managerActions,
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
					audioService.StopMusic()
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
	log.Debug().Msgf("Handle next alarm step at minute %d", r.ringingSinceMinutes)

	// Check for abort
	currentDelay := time.Duration(r.ringingSinceMinutes) * time.Minute
	abortDelay := r.lightDuration + 10*time.Minute
	if currentDelay > abortDelay {
		log.Warn().Dur("abort_delay", abortDelay).Msg("Ringer reached abort limit. Requesting manager to abort alarm.")
		r.managerActions <- ManagerActionAbortAlarm
		return
	}

	// Check for buzzer
	if currentDelay > r.lightDuration {
		log.Info().Dur("current_delay", currentDelay).Msg("Starting buzzer")
		r.physicalService.StartBuzzer()
	}

	// Check for music
	soundDelay := r.lightDuration - r.soundDuration
	if soundDelay.Truncate(time.Minute) == currentDelay.Truncate(time.Minute) {
		log.Info().Dur("current_delay", currentDelay).Msg("Starting alarm music")
		if err := r.audioService.PlayMusic(); err != nil {
			log.Error().Dur("current_delay", currentDelay).Msg("Failed to play music")
		}
	} else if currentDelay > soundDelay {
		log.Info().Dur("current_delay", currentDelay).Msg("Increasing alarm volume")
		if err := r.audioService.IncreaseVolume(); err != nil {
			log.Error().Dur("current_delay", currentDelay).Msg("Failed to increase alarm volume")
		}
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
