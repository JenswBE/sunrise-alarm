package entities

import (
	"time"

	"github.com/rs/zerolog"
)

var _ zerolog.LogObjectMarshaler = Planning{}

type Planning struct {
	NextSkipTime time.Time
	NextRingTime time.Time
}

func (p Planning) MarshalZerologObject(e *zerolog.Event) {
	e.Time("next_skip_time", p.NextSkipTime)
	e.Time("next_ring_time", p.NextRingTime)
}
