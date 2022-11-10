package entities

import (
	"time"

	"github.com/JenswBE/sunrise-alarm/src/entities"
)

type NextAlarmToRing struct {
	AlarmTime string `json:"alarm_time"`
}

func NextAlarmToRingResponseFromEntity(input entities.NextAlarm) NextAlarmToRing {
	return NextAlarmToRing{
		AlarmTime: input.AlarmTime.Format(time.RFC3339),
	}
}
