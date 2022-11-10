package usecases

import (
	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/JenswBE/sunrise-alarm/src/utils/pubsub"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *AlarmService) ListAlarms() ([]entities.Alarm, error) {
	log.Debug().Msg("Alarms Service: Listing alarms")
	return s.db.List()
}

func (s *AlarmService) GetAlarm(id uuid.UUID) (entities.Alarm, error) {
	log.Debug().Stringer("id", id).Msg("Alarms Service: Getting alarm")
	return s.db.Get(id)
}

// GetNextAlarmToRing returns the next alarm that will ring.
// If there are no future alarms, nil will be returned.
func (s *AlarmService) GetNextAlarmToRing() *entities.NextAlarm {
	return s.nextAlarmToRing
}

func (s *AlarmService) CreateAlarm(alarm entities.Alarm) (entities.Alarm, error) {
	log.Debug().Interface("alarm", alarm).Msg("Alarms Service: Creating alarm")
	newAlarm, err := s.db.Create(alarm)
	if err != nil {
		return entities.Alarm{}, err
	}
	s.publishAlarmsChanged("alarm_created")
	return newAlarm, nil
}

func (s *AlarmService) UpdateAlarm(alarm entities.Alarm) error {
	log.Debug().Interface("alarm", alarm).Msg("Alarms Service: Updating alarm")
	if err := s.db.Update(alarm); err != nil {
		return err
	}
	s.publishAlarmsChanged("alarm_updated")
	return nil
}

func (s *AlarmService) DeleteAlarm(id uuid.UUID) error {
	log.Debug().Stringer("id", id).Msg("Alarms Service: Deleting alarm")
	if err := s.db.Delete(id); err != nil {
		return err
	}
	s.publishAlarmsChanged("alarm_updated")
	return nil
}

func (s *AlarmService) publishAlarmsChanged(trigger string) {
	alarms, err := s.db.List()
	if err != nil {
		log.Error().Err(err).Str("trigger", trigger).Msg("Failed to list alarms to publish alarms_changed event")
		return
	}
	s.pubSub.Publish(&pubsub.EventAlarmsChanged{Alarms: alarms})
}
