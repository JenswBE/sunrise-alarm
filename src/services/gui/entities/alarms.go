package entities

import (
	"fmt"
	"strconv"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/samber/lo"
)

type AlarmsListTemplate struct {
	BaseData
	Alarms []entities.Alarm
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
		Days: lo.Map(e.Days, func(v time.Weekday, _ int) string { return strconv.Itoa(int(v)) }),
	}
}

func (e AlarmBody) ToEntity() (entities.Alarm, error) {
	a := entities.Alarm{
		Name: e.Name,
		Days: make([]time.Weekday, len(e.Days)),
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
		a.Days[i] = time.Weekday(daysInt)
	}
	return a, nil
}

func (e AlarmBody) HasWeekday(weekday time.Weekday) bool {
	weekdayString := strconv.Itoa(int(weekday))
	return lo.Contains(e.Days, weekdayString)
}

type AlarmsFormTemplate struct {
	BaseData
	AlarmBody       AlarmBody
	IsNew           bool
	Weekdays        []time.Weekday
	WeekdaysPresets map[string][]time.Weekday
}

func (t AlarmsFormTemplate) GetTemplateName() string {
	return "alarmsForm"
}
