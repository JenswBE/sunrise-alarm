package pubsub

import "github.com/rs/zerolog"

type Event interface {
	zerolog.LogObjectMarshaler
	GetTopic() string
}

type Channel chan Event

type PubSub interface {
	Subscribe(event Event, ch Channel)
	Publish(event Event)
}
