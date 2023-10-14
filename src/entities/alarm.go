package entities

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

const MaxNumberOfAlarms = 8

var _ zerolog.LogObjectMarshaler = Alarm{}

type Alarm struct {
	ID       uuid.UUID
	Enabled  bool
	Name     string
	Hour     uint8
	Minute   uint8
	Days     []ISOWeekday
	SkipNext bool
}

func (a Alarm) MarshalZerologObject(e *zerolog.Event) {
	e.Stringer("id", a.ID)
	e.Bool("enabled", a.Enabled)
	e.Str("name", a.Name)
	e.Str("time", a.TimeToString())
	days := lo.Map(a.Days, func(d ISOWeekday, _ int) string { return d.String() })
	e.Strs("days", days)
	e.Bool("skip_next", a.SkipNext)
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
	slices.SortStableFunc(alarms, func(a, b Alarm) int {
		switch {
		case a.Hour != b.Hour:
			return cmp.Compare(a.Hour, b.Hour)
		case a.Minute != b.Minute:
			return cmp.Compare(a.Minute, b.Minute)
		case len(a.Days) > 0 || len(b.Days) > 0:
			if len(a.Days) == 0 {
				return -1
			}
			if len(b.Days) == 0 {
				return +1
			}
			return cmp.Compare(a.Days[0], b.Days[0])
		default:
			return cmp.Compare(a.Name, b.Name)
		}
	})
}
