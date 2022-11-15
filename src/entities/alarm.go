package entities

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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
