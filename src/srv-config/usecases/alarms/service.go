package alarms

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-config/entities"
	"github.com/JenswBE/sunrise-alarm/src/srv-config/repositories"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var _ Usecase = &Service{}

const TopicPrefix = "sunrise_alarm/"

type Service struct {
	db         repositories.DB
	mqttClient repositories.MQTTClient
}

func NewService(db repositories.DB, mqttClient repositories.MQTTClient) *Service {
	return &Service{
		db:         db,
		mqttClient: mqttClient,
	}
}

func (s *Service) ListAlarms() ([]entities.Alarm, error) {
	log.Debug().Msg("Alarms Service: Listing alarms")
	return s.db.List()
}

func (s *Service) GetAlarm(id uuid.UUID) (entities.Alarm, error) {
	log.Debug().Stringer("id", id).Msg("Alarms Service: Getting alarm")
	return s.db.Get(id)
}

func (s *Service) CreateAlarm(alarm entities.Alarm) (entities.Alarm, error) {
	log.Debug().Interface("alarm", alarm).Msg("Alarms Service: Creating alarm")
	return s.db.Create(alarm)
}

func (s *Service) UpdateAlarm(alarm entities.Alarm) error {
	log.Debug().Interface("alarm", alarm).Msg("Alarms Service: Updating alarm")
	return s.db.Update(alarm)
}

func (s *Service) DeleteAlarm(id uuid.UUID) error {
	log.Debug().Stringer("id", id).Msg("Alarms Service: Deleting alarm")
	return s.db.Delete(id)
}
