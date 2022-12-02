package pubsub

import (
	"github.com/rs/zerolog"
)

var _ Event = EventButtonPressedShort{}

type EventButtonPressedShort struct{}

func (event EventButtonPressedShort) MarshalZerologObject(e *zerolog.Event) {
	e.Str("type", "ButtonPressedShort")
}

func (e EventButtonPressedShort) GetTopic() string {
	return "button_pressed_short"
}
