package gui

import (
	"testing"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/stretchr/testify/require"
)

func TestSortAlarms(t *testing.T) {
	testCases := map[string]struct {
		given    []entities.ISOWeekday
		expected string
	}{
		"nil": {
			given:    nil,
			expected: "_ _ _ _ _ _ _",
		},
		"empty slice": {
			given:    []entities.ISOWeekday{},
			expected: "_ _ _ _ _ _ _",
		},
		"odd days": {
			given: []entities.ISOWeekday{
				entities.ISOMonday,
				entities.ISOWednesday,
				entities.ISOFriday,
				entities.ISOSunday,
			},
			expected: "M _ W _ F _ S",
		},
		"even days": {
			given: []entities.ISOWeekday{
				entities.ISOTuesday,
				entities.ISOThursday,
				entities.ISOSaturday,
			},
			expected: "_ T _ T _ S _",
		},
		"all days": {
			given:    entities.ISOWeekdays(),
			expected: "M T W T F S S",
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			require.Equal(t, tc.expected, formatDays(tc.given))
		})
	}
}
