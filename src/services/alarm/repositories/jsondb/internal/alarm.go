package internal

import (
	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/google/uuid"
)

type Alarm struct {
	ID       uuid.UUID             `json:"id"`
	Enabled  bool                  `json:"enabled"`
	Name     string                `json:"name"`
	Hour     uint8                 `json:"hour"`
	Minute   uint8                 `json:"minute"`
	Days     []entities.ISOWeekday `json:"days"`
	SkipNext bool                  `json:"skip_next"`
}

func (a Alarm) ToGlobal() entities.Alarm {
	return entities.Alarm{
		ID:       a.ID,
		Enabled:  a.Enabled,
		Name:     a.Name,
		Hour:     a.Hour,
		Minute:   a.Minute,
		Days:     a.Days,
		SkipNext: a.SkipNext,
	}
}

func AlarmFromGlobal(e entities.Alarm) Alarm {
	return Alarm{
		ID:       e.ID,
		Enabled:  e.Enabled,
		Name:     e.Name,
		Hour:     e.Hour,
		Minute:   e.Minute,
		Days:     e.Days,
		SkipNext: e.SkipNext,
	}
}
