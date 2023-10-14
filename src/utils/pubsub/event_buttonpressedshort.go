package pubsub

import (
	"github.com/rs/zerolog"
)

var _ Event = EventButtonPressedShort{}

type EventButtonPressedShort struct{}

func (e EventButtonPressedShort) MarshalZerologObject(logEvent *zerolog.Event) {
	logEvent.Str("type", "ButtonPressedShort")
}

func (e EventButtonPressedShort) GetTopic() string {
	return "button_pressed_short"
}
