package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSortAlarms(t *testing.T) {
	singleAt7hNoName := Alarm{
		ID:   uuid.New(),
		Hour: 7,
	}
	singleAt7hA := Alarm{
		ID:   uuid.New(),
		Hour: 7,
		Name: "A",
	}
	singleAt7hB := Alarm{
		ID:   uuid.New(),
		Hour: 7,
		Name: "B",
	}
	singleAt7h30 := Alarm{
		ID:     uuid.New(),
		Hour:   7,
		Minute: 30,
	}
	mondayAt7h := Alarm{
		ID:   uuid.New(),
		Hour: 7,
		Days: []ISOWeekday{ISOMonday},
	}
	tuesdayAndWednesdayAt7h := Alarm{
		ID:   uuid.New(),
		Hour: 7,
		Days: []ISOWeekday{ISOTuesday, ISOWednesday},
	}
	sundayAt7h := Alarm{
		ID:   uuid.New(),
		Hour: 7,
		Days: []ISOWeekday{ISOSunday},
	}

	testCases := map[string]struct {
		given    []Alarm
		expected []Alarm
	}{
		"nil": {
			given:    nil,
			expected: nil,
		},
		"empty slice": {
			given:    []Alarm{},
			expected: []Alarm{},
		},
		"multiple alarms": {
			given: []Alarm{
				singleAt7h30,
				singleAt7hA,
				singleAt7hNoName,
				mondayAt7h,
				singleAt7hB,
				sundayAt7h,
				tuesdayAndWednesdayAt7h,
			},
			expected: []Alarm{
				singleAt7hNoName,
				singleAt7hA,
				singleAt7hB,
				mondayAt7h,
				tuesdayAndWednesdayAt7h,
				sundayAt7h,
				singleAt7h30,
			},
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			SortAlarms(tc.given)
			require.Equal(t, tc.expected, tc.given) // Alarms are sorted in place
		})
	}
}
