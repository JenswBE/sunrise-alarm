package usecases

import (
	"sync"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/services/audio"
	"github.com/JenswBE/sunrise-alarm/src/services/physical"
	"github.com/JenswBE/sunrise-alarm/src/utils/trigger"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

const (
	lightDelay    = 0 * time.Minute
	lightDuration = 10 * time.Minute
	lightDone     = lightDelay + lightDuration

	soundDelay    = 3 * time.Minute
	soundDuration = 7 * time.Minute
	soundDone     = soundDelay + soundDuration

	buzzerDelay    = 10 * time.Minute
	buzzerDuration = 5 * time.Minute
	buzzerDone     = buzzerDelay + buzzerDuration

	tickDuration = time.Minute
)

type Ringer struct {
	ringingMutex        sync.Mutex
	ringing             bool
	ringingChan         chan bool
	ringingSinceMinutes uint
	lightStatus         stepStatus
	soundStatus         stepStatus
	buzzerStatus        stepStatus
	abortAlarm          chan<- struct{}
	audioService        audio.Service
	physicalService     physical.Service
}

type stepStatus string

const (
	stepStatusAwaitingStart stepStatus = "AWAITING_START"
	stepStatusOngoing       stepStatus = "ONGOING"
	stepStatusDone          stepStatus = "DONE"
)

func NewRinger(physicalService physical.Service, audioService audio.Service, abortAlarm chan<- struct{}) *Ringer {
	// Init
	ringer := &Ringer{
		ringingChan:     make(chan bool, 1),
		abortAlarm:      abortAlarm,
		audioService:    audioService,
		physicalService: physicalService,
	}

	// Start event loop
	go ringer.eventLoop(abortAlarm)
	return ringer
}

func (r *Ringer) eventLoop(abortAlarm chan<- struct{}) {
	trigger := trigger.NewDelayedTrigger()
	for {
		select {
		case shouldRing := <-r.ringingChan:
			if shouldRing {
				r.ringingSinceMinutes = 0
				r.lightStatus = stepStatusAwaitingStart
				r.soundStatus = stepStatusAwaitingStart
				r.buzzerStatus = stepStatusAwaitingStart
				r.handleNextStep()
				trigger.Schedule(tickDuration)
			} else {
				if r.lightStatus == stepStatusOngoing {
					log.Debug().Msg("Ringer.handleNextStep: Stopping alarm light because alarm is stopped/aborted")
					r.physicalService.StopSunriseSimulation()
				}
				if r.soundStatus == stepStatusOngoing {
					log.Debug().Msg("Ringer.handleNextStep: Stopping alarm sound because alarm is stopped/aborted")
					if err := r.audioService.StopMusic(); err != nil {
						log.Error().Err(err).Msg("Ringer: Failed to stop music")
					}
				}
				if r.buzzerStatus == stepStatusOngoing {
					log.Debug().Msg("Ringer.handleNextStep: Stopping alarm buzzer because alarm is stopped/aborted")
					r.physicalService.StopBuzzer()
				}
			}
		case <-trigger.C:
			if !r.ringing {
				log.Debug().Msg("Ringer.loop: Received delayed trigger, but ringer is not ringing. Ignoring delayed trigger.")
				continue
			}
			r.ringingSinceMinutes++
			r.handleNextStep()
			trigger.Schedule(tickDuration)
		}
	}
}

func (r *Ringer) handleNextStep() {
	// Init
	currentDelay := time.Duration(r.ringingSinceMinutes) * time.Minute
	logger := log.With().Stringer("current_delay", currentDelay).Logger()
	logger.Debug().Msgf("Ringer.handleNextStep: Handle next alarm step at minute %d", r.ringingSinceMinutes)

	// Check for abort
	abortDelay := lo.Max([]time.Duration{lightDone, soundDone, buzzerDone})
	if currentDelay >= abortDelay {
		logger.Warn().Stringer("abort_delay", abortDelay).Msg("Ringer.handleNextStep: Ringer reached abort limit. Requesting manager to abort alarm.")
		r.abortAlarm <- struct{}{}
		return
	}

	// Check for light
	switch {
	case currentDelay >= lightDelay && r.lightStatus == stepStatusAwaitingStart:
		logger.Debug().Msg("Ringer.handleNextStep: Starting alarm light")
		r.physicalService.StartSunriseSimulation()
		r.lightStatus = stepStatusOngoing
	case currentDelay >= lightDone && r.lightStatus == stepStatusOngoing:
		logger.Debug().Msg("Ringer.handleNextStep: Stopping alarm light because done is reached")
		r.physicalService.StopSunriseSimulation()
		r.lightStatus = stepStatusDone
	}

	// Check for sound
	switch {
	case currentDelay >= soundDelay && r.soundStatus == stepStatusAwaitingStart:
		logger.Debug().Msg("Ringer.handleNextStep: Starting alarm sound")
		if err := r.audioService.PlayMusic(); err != nil {
			logger.Error().Msg("Ringer.handleNextStep: Failed to play music")
		}
		r.soundStatus = stepStatusOngoing
	case currentDelay < soundDone && r.soundStatus == stepStatusOngoing:
		logger.Debug().Msg("Ringer.handleNextStep: Increasing alarm sound volume")
		if err := r.audioService.IncreaseVolume(); err != nil {
			logger.Error().Msg("Ringer.handleNextStep: Failed to increase alarm sound volume")
		}
	case currentDelay >= soundDone && r.soundStatus == stepStatusOngoing:
		logger.Debug().Msg("Ringer.handleNextStep: Stopping alarm sound because done is reached")
		if err := r.audioService.StopMusic(); err != nil {
			log.Error().Err(err).Msg("Ringer: Failed to stop music")
		}
		r.soundStatus = stepStatusDone
	}

	// Check for buzzer
	switch {
	case currentDelay >= buzzerDelay && r.buzzerStatus == stepStatusAwaitingStart:
		logger.Debug().Msg("Ringer.handleNextStep: Starting alarm buzzer")
		r.physicalService.StartBuzzer()
		r.buzzerStatus = stepStatusOngoing
	case currentDelay >= buzzerDone && r.buzzerStatus == stepStatusOngoing:
		logger.Debug().Msg("Ringer.handleNextStep: Stopping alarm buzzer because done is reached")
		r.physicalService.StopBuzzer()
		r.buzzerStatus = stepStatusDone
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
		r.ringing = value
		r.ringingChan <- value
	}
}
