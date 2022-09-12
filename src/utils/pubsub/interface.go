package pubsub

type Event interface {
	GetTopic() string
}

type Channel chan Event

type PubSub interface {
	Subscribe(event Event, ch Channel)
	Publish(event Event)
}
