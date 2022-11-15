package trigger

import (
	"time"
)

type DelayedTrigger struct {
	C chan struct{}
}

func NewDelayedTrigger() *DelayedTrigger {
	return &DelayedTrigger{
		C: make(chan struct{}),
	}
}

// Schedule schedules a future trigger.
func (t *DelayedTrigger) Schedule(delay time.Duration) {
	go func() {
		time.Sleep(delay)
		t.Trigger()
	}()
}

// Trigger forces a manual trigger of the delayed trigger.
// Scheduled triggers are unharmed.
func (t *DelayedTrigger) Trigger() {
	t.C <- struct{}{}
}
