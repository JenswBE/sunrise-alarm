package pubsub

import "github.com/rs/zerolog"

var _ Event = EventButtonPressedLong{}

type EventButtonPressedLong struct{}

func (e EventButtonPressedLong) GetTopic() string {
	return "button_pressed_long"
}

func (event EventButtonPressedLong) MarshalZerologObject(e *zerolog.Event) {
	e.Str("type", "ButtonPressedLong")
}
