package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWeekdayAndISOWeekdayMapping(t *testing.T) {
	testCases := map[time.Weekday]ISOWeekday{
		time.Monday:    ISOMonday,
		time.Tuesday:   ISOTuesday,
		time.Wednesday: ISOWednesday,
		time.Thursday:  ISOThursday,
		time.Friday:    ISOFriday,
		time.Saturday:  ISOSaturday,
		time.Sunday:    ISOSunday,
	}

	for weekday, isoWeekday := range testCases {
		t.Run(weekday.String(), func(t *testing.T) {
			require.Equal(t, isoWeekday, NewISOWeekday(weekday))
			require.Equal(t, weekday, isoWeekday.ToWeekday())
		})
	}
}
