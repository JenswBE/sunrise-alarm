// Based on https://levelup.gitconnected.com/lets-write-a-simple-event-bus-in-go-79b9480d8997
package pubsub

import (
	"slices"
	"sync"

	"github.com/rs/zerolog/log"
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
	log.Debug().Object("event", event).Msgf("EventBus.Publish: New event published of type %T", event)
	if event == nil {
		log.Error().Msgf("EventBus.Publish: %T received with value nil. Pointer events are not supported and will be ignored.", event)
		return
	}
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
