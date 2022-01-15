package repositories

import "context"

type Button interface {
	IsPressed() bool
}

type MQTTClient interface {
	Publish(ctx context.Context, topic, payload string) error
}
