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

type EventAlarmsChanged struct {
	Alarms []entities.Alarm
}

func (e *EventAlarmsChanged) GetTopic() string {
	return "alarms_changed"
}

type EventNextAlarmsUpdated struct {
	ToRing     *entities.NextAlarm
	WithAction *entities.NextAlarm
}

func (e *EventNextAlarmsUpdated) GetTopic() string {
	return "next_alarms_updated"
}
