package presenter

import (
	"strconv"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/srv-config/api/openapi"
	"github.com/JenswBE/sunrise-alarm/src/srv-config/entities"
	"github.com/google/uuid"
)

func AlarmFromEntity(alarm entities.Alarm) openapi.Alarm {
	// Convert days
	days := make([]int32, len(alarm.Days))
	for i, day := range alarm.Days {
		// Convert to human indexed weekdays
		// Sunday = 0 and Monday = 1 => Monday = 1 and Sunday = 7
		if day == time.Sunday {
			days[i] = 7 // Convert Sunday 0 => 7
		} else {
			days[i] = int32(day)
		}
	}

	// Set basic fields
	return *openapi.NewAlarm(
		alarm.ID.String(),
		alarm.Enabled,
		alarm.Name,
		int32(alarm.Hour),
		int32(alarm.Minute),
		days,
		alarm.SkipNext,
	)
}

func AlarmToEntity(alarm openapi.Alarm) (entities.Alarm, error) {
	// Convert days
	openapiDays := alarm.GetDays()
	days := make([]time.Weekday, len(openapiDays))
	for i, day := range openapiDays {
		if day < 1 || day > 7 {
			return entities.Alarm{}, entities.NewError(400, openapi.ERRORCODE_INVALID_WEEKDAY, strconv.Itoa(int(day)), nil)
		}

		// Convert to machine indexed weekdays
		// Monday = 1 and Sunday = 7 => Sunday = 0 and Monday = 1
		if day == 7 {
			days[i] = time.Sunday
		} else {
			days[i] = time.Weekday(day)
		}
	}

	// Parse ID
	var id uuid.UUID
	if alarm.GetId() != "" {
		var err error
		id, err = ParseUUID(alarm.GetId())
		if err != nil {
			return entities.Alarm{}, err
		}
	}

	// Convert remaining fields
	return entities.Alarm{
		ID:       id,
		Enabled:  alarm.GetEnabled(),
		Name:     alarm.GetName(),
		Hour:     int8(alarm.GetHour()),
		Minute:   int8(alarm.GetMinute()),
		Days:     days,
		SkipNext: alarm.GetSkipNext(),
	}, nil
}
