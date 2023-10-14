package pubsub

import (
	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/rs/zerolog"
)

var _ Event = EventAlarmChanged{}

type EventAlarmChanged struct {
	Action AlarmChangedAction
	Alarm  entities.Alarm
}

func (e EventAlarmChanged) MarshalZerologObject(logEvent *zerolog.Event) {
	logEvent.Str("type", "AlarmChanged")
	logEvent.Stringer("action", e.Action)
	logEvent.Object("alarm", e.Alarm)
}

func (e EventAlarmChanged) GetTopic() string {
	return "alarm_changed"
}

type AlarmChangedAction string

const (
	AlarmChangedActionCreated AlarmChangedAction = "CREATED"
	AlarmChangedActionUpdated AlarmChangedAction = "UPDATED"
	AlarmChangedActionDeleted AlarmChangedAction = "DELETED"
)

func (action AlarmChangedAction) String() string {
	return string(action)
}
