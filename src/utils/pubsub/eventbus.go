// Based on https://levelup.gitconnected.com/lets-write-a-simple-event-bus-in-go-79b9480d8997
package pubsub

import (
	"sync"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

// EventBus is an in-memory pub/sub service based on channels.
// The zero value is ready to use.
type EventBus struct {
	subscribers map[string][]Channel
	mutex       sync.RWMutex
}

func (eb *EventBus) Subscribe(event Event, ch Channel) {
	topic := event.GetTopic()
	eb.mutex.Lock()
	if eb.subscribers == nil {
		eb.subscribers = make(map[string][]Channel, 1)
	}
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]Channel{}, ch)
	}
	eb.mutex.Unlock()
}

func (eb *EventBus) Publish(event Event) {
	log.Debug().Interface("event", event).Msgf("New event published of type %T", event)
	eb.mutex.RLock()
	if chans, found := eb.subscribers[event.GetTopic()]; found {
		go func(data Event, dataChannelSlices []Channel) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(event, slices.Clone(chans))
	}
	eb.mutex.RUnlock()
}
