package planner

import (
	"testing"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestCalculateNextTime(t *testing.T) {
	timezone := lo.Must(time.LoadLocation("Europe/Brussels"))
	now := time.Date(2020, time.September, 19, 14, 30, 0, 0, timezone) // Saturday
	today := []time.Weekday{now.Weekday()}
	weekdays := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday}
	testCases := map[string]struct {
		alarm                entities.Alarm
		now                  time.Time
		expectedNextSkipTime time.Time
		expectedNextRingTime time.Time
	}{
		"disabled": {
			alarm:                entities.Alarm{Enabled: false},
			now:                  now,
			expectedNextSkipTime: time.Time{},
			expectedNextRingTime: time.Time{},
		},
		"still today": {
			alarm:                entities.Alarm{Enabled: true, Hour: 16},
			now:                  now,
			expectedNextSkipTime: time.Time{},
			expectedNextRingTime: time.Date(2020, time.September, 19, 16, 0, 0, 0, timezone),
		},
		"for tomorrow": {
			alarm:                entities.Alarm{Enabled: true, Hour: 14},
			now:                  now,
			expectedNextSkipTime: time.Time{},
			expectedNextRingTime: time.Date(2020, time.September, 20, 14, 0, 0, 0, timezone),
		},
		"for today but skipped": {
			alarm:                entities.Alarm{Enabled: true, Hour: 16, SkipNext: true},
			now:                  now,
			expectedNextSkipTime: time.Date(2020, time.September, 19, 16, 0, 0, 0, timezone),
			expectedNextRingTime: time.Date(2020, time.September, 20, 16, 0, 0, 0, timezone),
		},
		"for tomorrow but skipped": {
			alarm:                entities.Alarm{Enabled: true, Hour: 14, SkipNext: true},
			now:                  now,
			expectedNextSkipTime: time.Date(2020, time.September, 20, 14, 0, 0, 0, timezone),
			expectedNextRingTime: time.Date(2020, time.September, 21, 14, 0, 0, 0, timezone),
		},
		"repeated and for today": {
			alarm:                entities.Alarm{Enabled: true, Hour: 16, Days: today},
			now:                  now,
			expectedNextSkipTime: time.Time{},
			expectedNextRingTime: time.Date(2020, time.September, 19, 16, 0, 0, 0, timezone),
		},
		"repeated but not for today": {
			alarm:                entities.Alarm{Enabled: true, Hour: 14, Days: today},
			now:                  now,
			expectedNextSkipTime: time.Time{},
			expectedNextRingTime: time.Date(2020, time.September, 26, 14, 0, 0, 0, timezone),
		},
		"repeated on weekdays time still to come": {
			alarm:                entities.Alarm{Enabled: true, Hour: 16, Days: weekdays},
			now:                  now,
			expectedNextSkipTime: time.Time{},
			expectedNextRingTime: time.Date(2020, time.September, 21, 16, 0, 0, 0, timezone),
		},
		"repeated on weekdays time past": {
			alarm:                entities.Alarm{Enabled: true, Hour: 14, Days: weekdays},
			now:                  now,
			expectedNextSkipTime: time.Time{},
			expectedNextRingTime: time.Date(2020, time.September, 21, 14, 0, 0, 0, timezone),
		},
		"repeated for today but skipped": {
			alarm:                entities.Alarm{Enabled: true, Hour: 16, Days: today, SkipNext: true},
			now:                  now,
			expectedNextSkipTime: time.Date(2020, time.September, 19, 16, 0, 0, 0, timezone),
			expectedNextRingTime: time.Date(2020, time.September, 26, 16, 0, 0, 0, timezone),
		},
		"repeated but not for today and skipped": {
			alarm:                entities.Alarm{Enabled: true, Hour: 14, Days: today, SkipNext: true},
			now:                  now,
			expectedNextSkipTime: time.Date(2020, time.September, 26, 14, 0, 0, 0, timezone),
			expectedNextRingTime: time.Date(2020, time.October, 3, 14, 0, 0, 0, timezone),
		},
		"repeated for weekdays but skipped and time still to come": {
			alarm:                entities.Alarm{Enabled: true, Hour: 16, Days: weekdays, SkipNext: true},
			now:                  now,
			expectedNextSkipTime: time.Date(2020, time.September, 21, 16, 0, 0, 0, timezone),
			expectedNextRingTime: time.Date(2020, time.September, 22, 16, 0, 0, 0, timezone),
		},
		"repeated for weekdays but skipped and time past": {
			alarm:                entities.Alarm{Enabled: true, Hour: 14, Days: weekdays, SkipNext: true},
			now:                  now,
			expectedNextSkipTime: time.Date(2020, time.September, 21, 14, 0, 0, 0, timezone),
			expectedNextRingTime: time.Date(2020, time.September, 22, 14, 0, 0, 0, timezone),
		},
		"handling DST start": {
			alarm:                entities.Alarm{Enabled: true, Hour: 14},
			now:                  time.Date(2022, time.March, 26, 14, 30, 0, 0, timezone), // Winter time
			expectedNextSkipTime: time.Time{},
			expectedNextRingTime: time.Date(2022, time.March, 27, 14, 0, 0, 0, timezone), // Summer time, note expected hour is same as alarm and not shifted due to DST start
		},
		"handling DST end": {
			alarm:                entities.Alarm{Enabled: true, Hour: 14},
			now:                  time.Date(2022, time.October, 29, 14, 30, 0, 0, timezone), // Summer time
			expectedNextSkipTime: time.Time{},
			expectedNextRingTime: time.Date(2022, time.October, 30, 14, 0, 0, 0, timezone), // Winter time, note expected hour is same as alarm and not shifted due to DST end
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			actual := CalculatePlanning(tc.alarm, tc.now)
			require.Equal(t, tc.expectedNextSkipTime, actual.NextSkipTime, "Next skip time should be equal")
			require.Equal(t, tc.expectedNextRingTime, actual.NextRingTime, "Next ring time should be equal")
		})
	}
}
