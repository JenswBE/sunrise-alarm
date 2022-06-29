package entities

import (
	"time"

	"github.com/google/uuid"
)

type Alarm struct {
	ID       uuid.UUID      `json:"id"`
	Enabled  bool           `json:"enabled"`
	Name     string         `json:"name"`
	Hour     int8           `json:"hour"`
	Minute   int8           `json:"minute"`
	Days     []time.Weekday `json:"days"`
	SkipNext bool           `json:"skip_next"`
}
