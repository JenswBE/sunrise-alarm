package planner

import (
	"time"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/samber/lo"
)

func CalculatePlanning(alarm entities.Alarm, now time.Time) entities.Planning {
	// Check enabled
	if !alarm.Enabled {
		return entities.Planning{}
	}

	// Calculate next alarm time
	nextRingTime := calculateNextTime(alarm, now)

	// Check if we need to skip
	if alarm.SkipNext {
		return entities.Planning{
			NextSkipTime: nextRingTime,
			NextRingTime: calculateNextTime(alarm, nextRingTime),
		}
	}

	// Return next alarm
	return entities.Planning{
		NextSkipTime: time.Time{},
		NextRingTime: nextRingTime,
	}
}

// calculateNextTime calculates the next time the alarm should trigger
func calculateNextTime(alarm entities.Alarm, after time.Time) time.Time {
	// Create possible next alarm date
	alarmTime := time.Date(after.Year(), after.Month(), after.Day(), int(alarm.Hour), int(alarm.Minute), 0, 0, after.Location())

	// Check if alarm is exactly at date
	if len(alarm.Days) == 0 || lo.Contains(alarm.Days, entities.NewISOWeekday(after.Weekday())) {
		if alarmTime.After(after) {
			// Alarm is at date and still to come
			return alarmTime
		}
	}

	// Next alarm is not on date => Calculate next weekday
	nextDay := calculateNextDay(alarm, entities.NewISOWeekday(after.Weekday()))
	dayDiff := calculateWeekdaysDiff(entities.NewISOWeekday(alarmTime.Weekday()), nextDay)
	if dayDiff == 0 {
		dayDiff = 7
	}
	return alarmTime.AddDate(0, 0, dayDiff)
}

// calculateNextDay calculates the next day the alarm should trigger
func calculateNextDay(alarm entities.Alarm, currentWeekday entities.ISOWeekday) entities.ISOWeekday {
	if len(alarm.Days) == 0 {
		// No days set => next day is tomorrow
		return (currentWeekday + 1) % 7
	}

	// Get next day => First day higher then current day
	nextDay, found := lo.Find(alarm.Days, func(d entities.ISOWeekday) bool { return d > currentWeekday })
	if found {
		return nextDay
	}

	// No higher day found => Return first from list
	return alarm.Days[0]
}

func calculateWeekdaysDiff(from, to entities.ISOWeekday) int {
	if to >= from {
		return int(to) - int(from)
	}
	return int(to) - int(from) + 7
}
