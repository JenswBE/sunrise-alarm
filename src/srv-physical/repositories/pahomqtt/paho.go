package pahomqtt

import (
	"context"
	"fmt"
	"net/url"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var _ repositories.MQTTClient = &PahoMQTT{}

type PahoMQTT struct {
	client *autopaho.ConnectionManager
}

type pahoLogger struct{}

func (pahoLogger) Println(v ...interface{}) {
	log.Print(v...)
}
func (pahoLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func NewPahoMQTT(ctx context.Context, mqttBroker string, mqttPort int) (*PahoMQTT, error) {
	// Create client config
	brokerURL, err := url.Parse(fmt.Sprintf(`mqtt://%s:%d`, mqttBroker, mqttPort))
	if err != nil {
		return nil, fmt.Errorf("failed to parse MQTT broker url: %w", err)
	}
	clientIDSuffix := uuid.New().String()[:8]
	config := autopaho.ClientConfig{
		BrokerUrls:     []*url.URL{brokerURL},
		KeepAlive:      10,
		OnConnectError: func(err error) { log.Error().Err(err).Msg("PahoMQTT: (Re)connecting to MQTT broker failed") },
		Debug:          pahoLogger{},
		ClientConfig: paho.ClientConfig{
			ClientID:      "srv-physical-" + clientIDSuffix,
			OnClientError: func(err error) { log.Error().Err(err).Msg("PahoMQTT: MQTT broker requested disconnect") },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					log.Error().Str("reason", d.Properties.ReasonString).Msg("PahoMQTT: MQTT broker requested disconnect")
				} else {
					log.Error().Uint8("reason_code", d.ReasonCode).Msg("PahoMQTT: MQTT broker requested disconnect")
				}
			},
		},
	}

	// Create client
	cm, err := autopaho.NewConnection(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new autopaho client: %w", err)
	}

	// Connect to broker
	err = cm.AwaitConnection(ctx)
	if err != nil {
		// Should only happen when context is cancelled
		return nil, fmt.Errorf("context cancelled while awaiting MQTT connection: %w", err)
	}

	// Connect successful
	return &PahoMQTT{client: cm}, nil
}

func (p *PahoMQTT) Publish(ctx context.Context, topic, payload string) error {
	log.Debug().Str("topic", topic).Str("payload", payload).Msg("PahoMQTT: Publishing message to broker")
	_, err := p.client.Publish(ctx, &paho.Publish{
		QoS:     2, // 0 - at most once; 1 - at least once; 2 - exactly once
		Topic:   topic,
		Payload: []byte(payload),
	})
	return err
}
