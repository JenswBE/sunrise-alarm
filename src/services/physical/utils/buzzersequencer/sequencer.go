package buzzersequencer

import (
	"time"

	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
)

type BuzzerSequencer struct {
	buzzer    repositories.Buzzer
	startStop chan bool
}

func NewBuzzerSequencer(buzzer repositories.Buzzer, fancySequence bool) *BuzzerSequencer {
	startStop := make(chan bool, 5)
	seq := &BuzzerSequencer{
		buzzer:    buzzer,
		startStop: startStop,
	}
	go seq.worker(startStop, fancySequence)
	return seq
}

func (s *BuzzerSequencer) Start() {
	s.startStop <- true
}

func (s *BuzzerSequencer) Stop() {
	s.startStop <- false
}

func (s *BuzzerSequencer) worker(startStop chan bool, fancySequence bool) {
	beepSequence := []time.Duration{
		500 * time.Millisecond, // On
		500 * time.Millisecond, // Off
	}
	if fancySequence {
		// Below sequence is a bit fancier,
		// but seems not always handled correctly.
		beepSequence = []time.Duration{
			100 * time.Millisecond, // On
			100 * time.Millisecond, // Off
			100 * time.Millisecond, // On
			1 * time.Second,        // Off
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
