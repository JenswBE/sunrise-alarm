package entities

import (
	"time"
)

type NextRingTime struct {
	RingTime string `json:"ring_time"`
}

func NextRingTimeResponseFromEntity(nextRingTime time.Time) NextRingTime {
	return NextRingTime{
		RingTime: nextRingTime.Format(time.RFC3339),
	}
}
