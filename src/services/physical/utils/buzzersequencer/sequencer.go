package buzzersequencer

import (
	"time"

	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
	"github.com/rs/zerolog/log"
)

type BuzzerSequencer struct {
	buzzer    repositories.Buzzer
	startStop chan bool
}

func NewBuzzerSequencer(buzzer repositories.Buzzer, beepSequence []time.Duration) *BuzzerSequencer {
	startStop := make(chan bool, 5)
	seq := &BuzzerSequencer{
		buzzer:    buzzer,
		startStop: startStop,
	}
	go seq.worker(startStop, beepSequence)
	return seq
}

func (s *BuzzerSequencer) Start() {
	s.startStop <- true
}

func (s *BuzzerSequencer) Stop() {
	s.startStop <- false
}

func (s *BuzzerSequencer) worker(startStop chan bool, beepSequence []time.Duration) {
	if len(beepSequence) == 0 {
		// Should not happen => Fallback sequence
		log.Error().Msg("Beep sequence not set. Falling back to safe value.")
		beepSequence = []time.Duration{
			200 * time.Millisecond, // On
			200 * time.Millisecond, // Off
		}
	}

	// Start worker loop
	for {
		for {
			// Await start signal
			if <-startStop {
				break
			}
		}

		// Beep sequence
		beepState := false
		beepStep := -1
	beepSequence:
		for {
			select {
			case start := <-startStop:
				// Quit sequence if stop received
				if !start {
					s.buzzer.Off() // Ensure buzzer is stopped
					break beepSequence
				}
			default:
				// Move to next step
				beepStep++
				if beepStep > len(beepSequence)-1 {
					beepStep = 0
				}

				// Toggle pin
				beepState = !beepState
				if beepState {
					s.buzzer.On()
				} else {
					s.buzzer.Off()
				}

				// Sleep until next step
				time.Sleep(beepSequence[beepStep])
			}
		}
	}
}
