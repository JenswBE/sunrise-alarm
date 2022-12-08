package tsl2591lightsensor

import (
	"time"

	tsl2591 "github.com/JenswBE/golang-tsl2591"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
	"github.com/rs/zerolog/log"
)

const Interval = 1 * time.Second

var _ repositories.LightSensor = &TSL2591LightSensor{}

type TSL2591LightSensor struct {
	sensor *tsl2591.TSL2591
}

func NewTSL2591LightSensor(i2cDevice string) (*TSL2591LightSensor, error) {
	sensor, err := tsl2591.NewTSL2591(&tsl2591.Opts{
		Bus:    i2cDevice,
		Gain:   tsl2591.GainMed,
		Timing: tsl2591.Integrationtime100MS,
	})
	if err != nil {
		log.Error().Err(err).Msg("TSL2591LightSensor.GetVisibleLight: Failed to open and configure sensor")
		return nil, err
	}
	return &TSL2591LightSensor{sensor: sensor}, nil
}

func (s *TSL2591LightSensor) GetVisibleLight() (uint32, error) {
	visible, err := s.sensor.Visible()
	if err != nil {
		log.Error().Err(err).Msg("TSL2591LightSensor.GetVisibleLight: Failed to get visible light from sensor")
		return 0, err
	}
	return visible, nil
}

func (s *TSL2591LightSensor) Close() error {
	err := s.sensor.Disable()
	if err != nil {
		log.Error().Err(err).Msg("TSL2591LightSensor.Close: Failed to close light sensor")
	}
	return err
}
