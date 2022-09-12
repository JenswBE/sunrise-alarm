package usecases

import (
	"time"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

func (s *AlarmService) calculateNextAlarms(alarms []entities.Alarm) (nextAlarmToRing, nextAlarmWithAction *entities.NextAlarm) {
	// Setup
	now := time.Now()

	// Get next alarms and remove disabled alarms
	nextAlarms := lo.FilterMap(alarms, func(a entities.Alarm, _ int) (entities.NextAlarm, bool) { return s.calculateNextAlarm(a, now) })
	if len(nextAlarms) == 0 {
		log.Debug().Msg("Planner.calculateNextAlarms: No next alarms found")
		return nil, nil
	}

	// Calculate next alarm to ring
	nextAlarmRing := lo.MinBy(nextAlarms, func(a entities.NextAlarm, min entities.NextAlarm) bool {
		return a.AlarmTime.Before(min.AlarmTime)
	})

	// Calculate next alarm with action
	nextAlarmAction := lo.MinBy(nextAlarms, func(a entities.NextAlarm, min entities.NextAlarm) bool {
		return a.NextActionTime.Before(min.NextActionTime)
	})
	return &nextAlarmRing, &nextAlarmAction
}

func (s *AlarmService) calculateNextAlarm(alarm entities.Alarm, now time.Time) (nextAlarm entities.NextAlarm, enabled bool) {
	// Check enabled
	if !alarm.Enabled {
		return
	}

	// Calculate next alarm time
	nextAlarmTime := calculateNextTime(alarm, now)

	// Check if we need to skip
	if alarm.SkipNext {
		return entities.NextAlarm{
			ID:             alarm.ID,
			AlarmTime:      calculateNextTime(alarm, nextAlarmTime),
			NextAction:     entities.NextActionSkip,
			NextActionTime: nextAlarmTime,
		}, true
	}

	// Check if we already ringed for this alarm
	if lastRing, found := s.lastRings[alarm.ID]; found {
		if lastRing.Equal(nextAlarmTime) || lastRing.After(nextAlarmTime) {
			// Alarm already ringed => Skip one iteration
			nextAlarmTime = calculateNextTime(alarm, nextAlarmTime)
		}
	}

	// Return next alarm
	return entities.NextAlarm{
		ID:             alarm.ID,
		AlarmTime:      nextAlarmTime,
		NextAction:     entities.NextActionRing,
		NextActionTime: nextAlarmTime.Add(-s.alarmLightDuration),
	}, true
}

// calculateNextTime calculates the next time the alarm should trigger
func calculateNextTime(alarm entities.Alarm, after time.Time) time.Time {
	// Create possible next alarm date
	alarmTime := time.Date(after.Year(), after.Month(), after.Day(), int(alarm.Hour), int(alarm.Minute), 0, 0, after.Location())

	// Check if alarm is exactly at date
	if len(alarm.Days) == 0 || lo.Contains(alarm.Days, after.Weekday()) {
		if alarmTime.After(after) {
			// Alarm is at date and still to come
			return alarmTime
		}
	}

	// Next alarm is not on date => Calculate next weekday
	nextDay := calculateNextDay(alarm, after.Weekday())
	dayDiff := calculateWeekdaysDiff(alarmTime.Weekday(), nextDay)
	if dayDiff == 0 {
		dayDiff = 7
	}
	return alarmTime.AddDate(0, 0, dayDiff)
}

// calculateNextDay calculates the next day the alarm should trigger
func calculateNextDay(alarm entities.Alarm, currentWeekday time.Weekday) time.Weekday {
	if len(alarm.Days) == 0 {
		// No days set => next day is tomorrow
		return (currentWeekday + 1) % 7
	}

	// Get next day => First day higher then current day
	nextDay, found := lo.Find(alarm.Days, func(d time.Weekday) bool { return d > currentWeekday })
	if found {
		return nextDay
	}

	// No higher day found => Return first from list
	return alarm.Days[0]
}

func calculateWeekdaysDiff(from, to time.Weekday) int {
	if to >= from {
		return int(to) - int(from)
	}
	return int(to) - int(from) + 7
}
