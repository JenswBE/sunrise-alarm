package debouncedbutton

import (
	"time"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories"
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

type DebouncedButton struct {
	button              repositories.Button
	notifyChannel       chan ButtonPress
	firstPressTimestamp time.Time
	isLongPress         bool
}

func NewDebouncedButton(button repositories.Button, notifyChannel chan ButtonPress) *DebouncedButton {
	dButton := &DebouncedButton{
		button:        button,
		notifyChannel: notifyChannel,
		isLongPress:   false,
	}
	go dButton.watcher()
	return dButton
}

func (b *DebouncedButton) watcher() {
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

func (b *DebouncedButton) handlePress() {
	// Ignore if we are in a long press
	if b.isLongPress {
		return
	}

	// Check if we reached a long press
	if !b.firstPressTimestamp.IsZero() && time.Since(b.firstPressTimestamp) > LongPressDuration {
		// Long press reached
		b.notifyChannel <- ButtonPressLong
		b.isLongPress = true
		return
	}

	// Save timestamp if first press
	if b.firstPressTimestamp.IsZero() {
		b.firstPressTimestamp = time.Now()
	}
}

func (b *DebouncedButton) handleRelease() {
	// Send event on short press
	if !b.isLongPress {
		b.notifyChannel <- ButtonPressShort
	}

	// Reset button
	b.isLongPress = false
	b.firstPressTimestamp = time.Time{}
}
