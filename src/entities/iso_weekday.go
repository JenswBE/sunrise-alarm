package entities

import "time"

// A Weekday specifies a day of the week accordingly to ISO (Monday = 0, ...).
// Contrary to the build-in ISOWeekday, this one starts on Monday instead of Sunday.
type ISOWeekday int

const (
	ISOMonday ISOWeekday = iota
	ISOTuesday
	ISOWednesday
	ISOThursday
	ISOFriday
	ISOSaturday
	ISOSunday
)

func ISOWeekdays() []ISOWeekday {
	return []ISOWeekday{
		ISOMonday,
		ISOTuesday,
		ISOWednesday,
		ISOThursday,
		ISOFriday,
		ISOSaturday,
		ISOSunday,
	}
}

func NewISOWeekday(w time.Weekday) ISOWeekday {
	if w == time.Sunday {
		// We consider Sunday the last day of the week
		return ISOSunday
	}
	// Shift weekdays 1 day down
	return ISOWeekday(w - 1)
}

func (w ISOWeekday) ToWeekday() time.Weekday {
	// Check if Sunday, which we consider last day of the week.
	if w == ISOSunday {
		return time.Sunday
	}
	// Shift weekdays 1 day up
	return time.Weekday(w + 1)
}

func (w ISOWeekday) String() string {
	return w.ToWeekday().String()
}
