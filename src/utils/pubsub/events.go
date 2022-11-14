package pubsub

import "github.com/JenswBE/sunrise-alarm/src/entities"

type EventButtonPressedShort struct{}

func (e *EventButtonPressedShort) GetTopic() string {
	return "button_pressed_short"
}

type EventButtonPressedLong struct{}

func (e *EventButtonPressedLong) GetTopic() string {
	return "button_pressed_long"
}

type EventAlarmChanged struct {
	Action AlarmChangedAction
	Alarm  entities.Alarm
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

func (e *EventAlarmChanged) GetTopic() string {
	return "alarm_changed"
}
