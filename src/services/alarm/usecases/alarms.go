package usecases

import (
	"fmt"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/JenswBE/sunrise-alarm/src/utils/pubsub"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

func (s *AlarmService) ListAlarms() ([]entities.Alarm, error) {
	log.Debug().Msg("Alarms Service: Listing alarms")
	return s.db.List()
}

func (s *AlarmService) GetAlarm(id uuid.UUID) (entities.Alarm, error) {
	log.Debug().Stringer("id", id).Msg("Alarms Service: Getting alarm")
	return s.db.Get(id)
}

// GetNextRingTime returns the next time an alarm will ring.
// If there are no future alarms, a zero time.Time will be returned.
func (s *AlarmService) GetNextRingTime() time.Time {
	plannings := lo.Values(s.planningsByAlarmID)
	firstPlanning := lo.MinBy(plannings, func(a, b entities.Planning) bool {
		return a.NextRingTime.Before(b.NextRingTime)
	})
	return firstPlanning.NextRingTime
}

func (s *AlarmService) CreateAlarm(alarm entities.Alarm) (entities.Alarm, error) {
	log.Debug().Object("alarm", alarm).Msg("Alarms Service: Creating alarm")
	existingAlarms, err := s.db.List()
	if err != nil {
		return entities.Alarm{}, err
	}
	if len(existingAlarms) >= entities.MaxNumberOfAlarms {
		return entities.Alarm{}, fmt.Errorf("maximum number of %d alarms reached to prevent display issues", entities.MaxNumberOfAlarms)
	}

	newAlarm, err := s.db.Create(alarm)
	if err != nil {
		return entities.Alarm{}, err
	}
	s.pubSub.Publish(pubsub.EventAlarmChanged{
		Action: pubsub.AlarmChangedActionCreated,
		Alarm:  newAlarm,
	})
	return newAlarm, nil
}

func (s *AlarmService) UpdateAlarm(alarm entities.Alarm) error {
	log.Debug().Object("alarm", alarm).Msg("Alarms Service: Updating alarm")
	if err := s.db.Update(alarm); err != nil {
		return err
	}
	s.pubSub.Publish(pubsub.EventAlarmChanged{
		Action: pubsub.AlarmChangedActionUpdated,
		Alarm:  alarm,
	})
	return nil
}

func (s *AlarmService) DeleteAlarm(id uuid.UUID) error {
	log.Debug().Stringer("id", id).Msg("Alarms Service: Deleting alarm")
	if err := s.db.Delete(id); err != nil {
		return err
	}
	s.pubSub.Publish(pubsub.EventAlarmChanged{
		Action: pubsub.AlarmChangedActionDeleted,
		Alarm:  entities.Alarm{ID: id},
	})
	return nil
}
