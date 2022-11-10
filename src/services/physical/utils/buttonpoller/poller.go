package buttonpoller

import (
	"time"

	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
	"github.com/rs/zerolog/log"
)

const (
	PollInterval      = 10 * time.Millisecond
	DebounceDuration  = 20 * time.Millisecond
	LongPressDuration = time.Second
)

type ButtonPress int

const (
	ButtonPressShort ButtonPress = iota
	ButtonPressLong
)

// ButtonPoller polls a button and generates short and long press events.
// The button is debounced as well.
type ButtonPoller struct {
	button              repositories.Button
	notifyChannel       chan ButtonPress
	firstPressTimestamp time.Time
	isLongPress         bool
}

func NewButtonPoller(button repositories.Button, notifyChannel chan ButtonPress) *ButtonPoller {
	dButton := &ButtonPoller{
		button:        button,
		notifyChannel: notifyChannel,
		isLongPress:   false,
	}
	go dButton.watcher()
	return dButton
}

func (b *ButtonPoller) watcher() {
	statePressed := b.button.IsPressed()
	var isPressed bool
	var isPressedChangedTS time.Time
	for {
		// Don't take all CPU
		time.Sleep(PollInterval)

		isPressed = b.button.IsPressed()
		if statePressed && isPressed {
			// Reconfirm button is still pressed
			b.handlePress()
			continue
		}

		if statePressed != isPressed {
			if isPressedChangedTS.IsZero() {
				// Start debounce timer
				isPressedChangedTS = time.Now()
			} else if time.Since(isPressedChangedTS) > DebounceDuration {
				// State changed
				statePressed = isPressed
				if statePressed {
					b.handlePress()
				} else {
					b.handleRelease()
				}
			}
		} else if !isPressedChangedTS.IsZero() {
			// Reset debounce
			isPressedChangedTS = time.Time{}
		}
	}
}

func (b *ButtonPoller) handlePress() {
	// Ignore if we are in a long press
	if b.isLongPress {
		return
	}

	// Check if we reached a long press
	if !b.firstPressTimestamp.IsZero() && time.Since(b.firstPressTimestamp) > LongPressDuration {
		// Long press reached
		log.Debug().Msg("ButtonPoller: Long press detected")
		b.notifyChannel <- ButtonPressLong
		b.isLongPress = true
		return
	}

	// Save timestamp if first press
	if b.firstPressTimestamp.IsZero() {
		b.firstPressTimestamp = time.Now()
	}
}

func (b *ButtonPoller) handleRelease() {
	// Send event on short press
	if !b.isLongPress {
		log.Debug().Msg("ButtonPoller: Short press detected")
		b.notifyChannel <- ButtonPressShort
	}

	// Reset button
	log.Debug().Msg("ButtonPoller: Reset button because of release")
	b.isLongPress = false
	b.firstPressTimestamp = time.Time{}
}
