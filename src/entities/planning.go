package entities

import (
	"time"
)

type Planning struct {
	NextSkipTime time.Time
	NextRingTime time.Time
}
