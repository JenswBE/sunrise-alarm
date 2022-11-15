package entities

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

const MaxNumberOfAlarms = 8

type Alarm struct {
	ID       uuid.UUID
	Enabled  bool
	Name     string
	Hour     uint8
	Minute   uint8
	Days     []time.Weekday
	SkipNext bool
}

func (a Alarm) TimeToString() string {
	return fmt.Sprintf("%02d:%02d", a.Hour, a.Minute)
}

func (a *Alarm) SetTimeFromString(time string) error {
	// Split parts
	timeParts := strings.Split(time, ":")
	if len(timeParts) != 2 {
		return fmt.Errorf("expected 2 time parts, but received %d: %s", len(timeParts), time)
	}
	hourString := strings.TrimPrefix(timeParts[0], "0")
	minuteString := strings.TrimPrefix(timeParts[1], "0")

	// Parse parts
	hour, err := strconv.ParseUint(hourString, 10, 8)
	if err != nil {
		return fmt.Errorf(`failed to parse hour component of time "%s": %w`, time, err)
	}
	minute, err := strconv.ParseUint(minuteString, 10, 8)
	if err != nil {
		return fmt.Errorf(`failed to parse minute component of time "%s": %w`, time, err)
	}

	// Update alarm
	a.Hour = uint8(hour)
	a.Minute = uint8(minute)
	return nil
}

func SortAlarms(alarms []Alarm) {
	slices.SortFunc(alarms, func(a, b Alarm) bool {
		switch {
		case a.Hour != b.Hour:
			return a.Hour < b.Hour
		case a.Minute != b.Minute:
			return a.Minute < b.Minute
		case len(a.Days) > 0 || len(b.Days) > 0:
			if len(a.Days) == 0 {
				return true
			}
			if len(b.Days) == 0 {
				return false
			}
			firstDayA := a.Days[0]
			if firstDayA == 0 {
				// We consider Sunday last day of the week
				firstDayA = 6
			}
			firstDayB := b.Days[0]
			if firstDayB == 0 {
				// We consider Sunday last day of the week
				firstDayB = 6
			}
			return firstDayA < firstDayB
		default:
			return a.Name < b.Name
		}
	})
}
