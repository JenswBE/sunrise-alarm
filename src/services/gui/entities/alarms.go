package entities

import (
	"fmt"
	"strconv"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/samber/lo"
)

type AlarmsListTemplate struct {
	BaseData
	AlarmsByStatus           []AlarmsWithStatus
	MaxNumberOfAlarmsReached bool
}

type AlarmsWithStatus struct {
	Alarms []entities.Alarm
	Status AlarmStatus
}

type AlarmStatus string

const (
	AlarmStatusEnabled  AlarmStatus = "ENABLED"
	AlarmStatusDisabled AlarmStatus = "DISABLED"
)

func (s AlarmStatus) String() string {
	return string(s)
}

func (s AlarmStatus) IsEnabled() bool {
	return s == AlarmStatusEnabled
}

func (s AlarmStatus) IsDisabled() bool {
	return s == AlarmStatusDisabled
}

func (t AlarmsListTemplate) GetTemplateName() string {
	return "alarmsList"
}

type AlarmBody struct {
	Name string   `form:"name"`
	Time string   `form:"time"`
	Days []string `form:"days"`
}

func AlarmBodyFromEntity(e entities.Alarm) AlarmBody {
	return AlarmBody{
		Name: e.Name,
		Time: e.TimeToString(),
		Days: lo.Map(e.Days, func(v entities.ISOWeekday, _ int) string { return strconv.Itoa(int(v)) }),
	}
}

func (e AlarmBody) ToEntity() (entities.Alarm, error) {
	a := entities.Alarm{
		Name: e.Name,
		Days: make([]entities.ISOWeekday, len(e.Days)),
	}
	if err := a.SetTimeFromString(e.Time); err != nil {
		return a, fmt.Errorf("failed to set time on Alarm entity from AlarmBody: %w", err)
	}
	for i, dayString := range e.Days {
		daysInt, err := strconv.Atoi(dayString)
		if err != nil {
			return a, fmt.Errorf("failed to parse days from AlarmBody %v: %w", e.Days, err)
		}
		if daysInt < 0 || daysInt > 6 {
			return a, fmt.Errorf("received invalid day %d. Days should be between 0 and 6 (both including)", daysInt)
		}
		a.Days[i] = entities.ISOWeekday(daysInt)
	}
	return a, nil
}

func (e AlarmBody) HasWeekday(weekday entities.ISOWeekday) bool {
	weekdayString := strconv.Itoa(int(weekday))
	return lo.Contains(e.Days, weekdayString)
}

type AlarmsFormTemplate struct {
	BaseData
	AlarmBody       AlarmBody
	IsNew           bool
	Weekdays        []entities.ISOWeekday
	WeekdaysPresets map[string][]entities.ISOWeekday
}

func (t AlarmsFormTemplate) GetTemplateName() string {
	return "alarmsForm"
}
