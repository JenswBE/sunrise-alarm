package entities

import (
	"time"

	"github.com/google/uuid"
)

type NextAlarm struct {
	ID             uuid.UUID
	AlarmTime      time.Time
	NextAction     NextAction
	NextActionTime time.Time
}

type NextAction string

const (
	NextActionRing NextAction = "RING"
	NextActionSkip NextAction = "SKIP"
)
