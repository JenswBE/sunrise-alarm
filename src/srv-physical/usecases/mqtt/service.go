package mqtt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories"
	"github.com/rs/zerolog/log"
)

var _ Usecase = &Service{}

const TopicPrefix = "sunrise_alarm/"

type Service struct {
	client repositories.MQTTClient
}

func NewService(mqttClient repositories.MQTTClient) *Service {
	return &Service{client: mqttClient}
}

func (s *Service) PublishButtonPressed(ctx context.Context) error {
	log.Debug().Msg("MQTT Service: Publishing button_pressed message to broker")
	return s.publishEmpty(ctx, "button_pressed")
}

func (s *Service) PublishButtonLongPressed(ctx context.Context) error {
	log.Debug().Msg("MQTT Service: Publishing button_long_pressed message to broker")
	return s.publishEmpty(ctx, "button_long_pressed")
}

func (s *Service) PublishTempHumidUpdated(ctx context.Context, tempHumid entities.TempHumidReading) error {
	log.Debug().Interface("temp_humid", tempHumid).Msg("MQTT Service: Publishing temp_humid_updated message to broker")
	return s.publishJSON(ctx, "temp_humid_updated", tempHumid)
}

func (s *Service) publishEmpty(ctx context.Context, topic string) error {
	return s.client.Publish(ctx, TopicPrefix+topic, "")
}

func (s *Service) publishJSON(ctx context.Context, topic string, msg interface{}) error {
	jsonMsg := bytes.Buffer{}
	err := json.NewEncoder(&jsonMsg).Encode(msg)
	if err != nil {
		return fmt.Errorf("failed to encode MQTT message into JSON: %w", err)
	}
	return s.client.Publish(ctx, TopicPrefix+topic, jsonMsg.String())
}
