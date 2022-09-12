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

func (t *DelayedTrigger) Schedule(delay time.Duration) {
	go func() {
		time.Sleep(delay)
		t.C <- struct{}{}
	}()
}
