package usecases

import (
	"testing"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestCalculateNextTime(t *testing.T) {
	timezone := lo.Must(time.LoadLocation("Europe/Brussels"))
	now := time.Date(2020, time.September, 19, 14, 30, 0, 0, timezone) // Saturday
	today := []time.Weekday{now.Weekday()}
	weekdays := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday}
	testCases := map[string]struct {
		alarm                  entities.Alarm
		now                    time.Time
		lastRinged             time.Time
		expectedAlarmTime      time.Time
		expectedNextAction     entities.NextAction
		expectedNextActionTime time.Time
	}{
		"disabled": {
			alarm: entities.Alarm{Enabled: false},
			now:   now,
		},
		"still today": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 16},
			now:                    now,
			expectedAlarmTime:      time.Date(2020, time.September, 19, 16, 0, 0, 0, timezone),
			expectedNextAction:     entities.NextActionRing,
			expectedNextActionTime: time.Date(2020, time.September, 19, 15, 53, 0, 0, timezone),
		},
		"for tomorrow": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 14},
			now:                    now,
			expectedAlarmTime:      time.Date(2020, time.September, 20, 14, 0, 0, 0, timezone),
			expectedNextAction:     entities.NextActionRing,
			expectedNextActionTime: time.Date(2020, time.September, 20, 13, 53, 0, 0, timezone),
		},
		"for today but skipped": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 16, SkipNext: true},
			now:                    now,
			expectedAlarmTime:      time.Date(2020, time.September, 20, 16, 0, 0, 0, timezone),
			expectedNextAction:     entities.NextActionSkip,
			expectedNextActionTime: time.Date(2020, time.September, 19, 16, 0, 0, 0, timezone),
		},
		"for tomorrow but skipped": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 14, SkipNext: true},
			now:                    now,
			expectedAlarmTime:      time.Date(2020, time.September, 21, 14, 0, 0, 0, timezone),
			expectedNextAction:     entities.NextActionSkip,
			expectedNextActionTime: time.Date(2020, time.September, 20, 14, 0, 0, 0, timezone),
		},
		"repeated and for today": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 16, Days: today},
			now:                    now,
			expectedAlarmTime:      time.Date(2020, time.September, 19, 16, 0, 0, 0, timezone),
			expectedNextAction:     entities.NextActionRing,
			expectedNextActionTime: time.Date(2020, time.September, 19, 15, 53, 0, 0, timezone),
		},
		"repeated and for today but already ringed": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 16, Days: today},
			now:                    now,
			lastRinged:             time.Date(2020, time.September, 19, 16, 0, 0, 0, timezone),
			expectedAlarmTime:      time.Date(2020, time.September, 26, 16, 0, 0, 0, timezone),
			expectedNextAction:     entities.NextActionRing,
			expectedNextActionTime: time.Date(2020, time.September, 26, 15, 53, 0, 0, timezone),
		},
		"repeated but not for today": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 14, Days: today},
			now:                    now,
			expectedAlarmTime:      time.Date(2020, time.September, 26, 14, 0, 0, 0, timezone),
			expectedNextAction:     entities.NextActionRing,
			expectedNextActionTime: time.Date(2020, time.September, 26, 13, 53, 0, 0, timezone),
		},
		"repeated on weekdays time still to come": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 16, Days: weekdays},
			now:                    now,
			expectedAlarmTime:      time.Date(2020, time.September, 21, 16, 0, 0, 0, timezone),
			expectedNextAction:     entities.NextActionRing,
			expectedNextActionTime: time.Date(2020, time.September, 21, 15, 53, 0, 0, timezone),
		},
		"repeated on weekdays time past": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 14, Days: weekdays},
			now:                    now,
			expectedAlarmTime:      time.Date(2020, time.September, 21, 14, 0, 0, 0, timezone),
			expectedNextAction:     entities.NextActionRing,
			expectedNextActionTime: time.Date(2020, time.September, 21, 13, 53, 0, 0, timezone),
		},
		"repeated for today but skipped": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 16, Days: today, SkipNext: true},
			now:                    now,
			expectedAlarmTime:      time.Date(2020, time.September, 26, 16, 0, 0, 0, timezone),
			expectedNextAction:     entities.NextActionSkip,
			expectedNextActionTime: time.Date(2020, time.September, 19, 16, 0, 0, 0, timezone),
		},
		"repeated but not for today and skipped": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 14, Days: today, SkipNext: true},
			now:                    now,
			expectedAlarmTime:      time.Date(2020, time.October, 3, 14, 0, 0, 0, timezone),
			expectedNextAction:     entities.NextActionSkip,
			expectedNextActionTime: time.Date(2020, time.September, 26, 14, 0, 0, 0, timezone),
		},
		"repeated for weekdays but skipped and time still to come": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 16, Days: weekdays, SkipNext: true},
			now:                    now,
			expectedAlarmTime:      time.Date(2020, time.September, 22, 16, 0, 0, 0, timezone),
			expectedNextAction:     entities.NextActionSkip,
			expectedNextActionTime: time.Date(2020, time.September, 21, 16, 0, 0, 0, timezone),
		},
		"repeated for weekdays but skipped and time past": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 14, Days: weekdays, SkipNext: true},
			now:                    now,
			expectedAlarmTime:      time.Date(2020, time.September, 22, 14, 0, 0, 0, timezone),
			expectedNextAction:     entities.NextActionSkip,
			expectedNextActionTime: time.Date(2020, time.September, 21, 14, 0, 0, 0, timezone),
		},
		"handling DST start": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 14},
			now:                    time.Date(2022, time.March, 26, 14, 30, 0, 0, timezone), // Winter time
			expectedAlarmTime:      time.Date(2022, time.March, 27, 14, 0, 0, 0, timezone),  // Summer time, note expected hour is same as alarm and not shifted due to DST start
			expectedNextAction:     entities.NextActionRing,
			expectedNextActionTime: time.Date(2022, time.March, 27, 13, 53, 0, 0, timezone),
		},
		"handling DST end": {
			alarm:                  entities.Alarm{Enabled: true, Hour: 14},
			now:                    time.Date(2022, time.October, 29, 14, 30, 0, 0, timezone), // Summer time
			expectedAlarmTime:      time.Date(2022, time.October, 30, 14, 0, 0, 0, timezone),  // Winter time, note expected hour is same as alarm and not shifted due to DST end
			expectedNextAction:     entities.NextActionRing,
			expectedNextActionTime: time.Date(2022, time.October, 30, 13, 53, 0, 0, timezone),
		},
	}

	for desc, tc := range testCases {
		s := AlarmService{alarmLightDuration: 7 * time.Minute}
		t.Run(desc, func(t *testing.T) {
			if !tc.lastRinged.IsZero() {
				s.lastRings = map[uuid.UUID]time.Time{tc.alarm.ID: tc.lastRinged}
			}
			actual, enabled := s.calculateNextAlarm(tc.alarm, tc.now)
			if tc.expectedAlarmTime.IsZero() {
				require.False(t, enabled, "Alarm is expected to be disabled")
			} else {
				require.True(t, enabled, "Alarm is expected to be enabled")
				require.Equal(t, tc.expectedAlarmTime, actual.AlarmTime, "Alarm time should be equal")
				require.Equal(t, tc.expectedNextAction, actual.NextAction, "Next action should be equal")
				require.Equal(t, tc.expectedNextActionTime, actual.NextActionTime, "Next action time should be equal")
			}
		})
	}
}
