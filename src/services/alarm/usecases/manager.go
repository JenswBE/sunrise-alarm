package usecases

import (
	"time"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	physicalEntities "github.com/JenswBE/sunrise-alarm/src/services/physical/entities"
	"github.com/JenswBE/sunrise-alarm/src/utils/pubsub"
	"github.com/JenswBE/sunrise-alarm/src/utils/trigger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type ManagerAction string

const (
	ManagerActionAbortAlarm     ManagerAction = "ABORT_ALARM"
	ManagerActionUpdateSchedule ManagerAction = "UPDATE_SCHEDULE"
)

func (action ManagerAction) String() string {
	return string(action)
}

func (s *AlarmService) startManager(managerActions chan ManagerAction) error {
	// Calculate initial next alarms
	alarms, err := s.ListAlarms()
	if err != nil {
		return err
	}
	s.nextAlarmToRing, s.nextAlarmWithAction = s.calculateNextAlarms(alarms)
	log.Debug().Interface("next_alarm_to_ring", s.nextAlarmToRing).Interface("next_alarm_with_action", s.nextAlarmWithAction).Msg("Manager: Initial next alarms calculated")

	// Setup initial delay
	initialDelay := s.durationUntilNextAction(time.Second)
	loggerWithFormattedDelay(initialDelay).Info().Msg("Initial delay set")

	go func(initialDelay time.Duration) {
		trigger := trigger.NewDelayedTrigger()
		trigger.Schedule(initialDelay)
		events := make(chan pubsub.Event, 1)
		s.pubSub.Subscribe(&pubsub.EventAlarmsChanged{}, events)
		s.pubSub.Subscribe(&pubsub.EventButtonPressedShort{}, events)
		s.pubSub.Subscribe(&pubsub.EventButtonPressedLong{}, events)
		for {
			select {
			case event := <-events:
				switch e := event.(type) {
				case *pubsub.EventAlarmsChanged:
					log.Debug().Msg("Manager: EventAlarmsChanged event received, requesting planner to calculate next alarms...")
					s.nextAlarmToRing, s.nextAlarmWithAction = s.calculateNextAlarms(e.Alarms)
					log.Debug().Interface("next_alarm_to_ring", s.nextAlarmToRing).Interface("next_alarm_with_action", s.nextAlarmWithAction).Msg("Manager: Next alarms calculated")
					managerActions <- ManagerActionUpdateSchedule
				case *pubsub.EventButtonPressedShort:
					log.Debug().Msg("Manager: EventButtonPressedShort event received, handling button press...")
					s.handleButtonPressed()
				case *pubsub.EventButtonPressedLong:
					log.Debug().Msg("Manager: EventButtonPressedLong event received, handling button press...")
					s.handleButtonLongPressed()
				}
			case action := <-managerActions:
				newDelay := s.handleManagerAction(action)
				if newDelay != nil {
					loggerWithFormattedDelay(*newDelay).Debug().Msg("Manager: Set new delayed trigger")
					trigger.Schedule(*newDelay)
				}
			case <-trigger.C:
				// Delay will be set by handle_action (through UpdateSchedule)
				log.Debug().Msg("Manager: Received delayed trigger, calling handleNextAction...")
				s.handleNextAction()
			}
		}
	}(initialDelay)

	return nil
}

func loggerWithFormattedDelay(delay time.Duration) *zerolog.Logger {
	logger := log.With().
		Stringer("delay", delay).
		Float64("delay_days", float64(delay)/float64(24*time.Hour)).
		Logger()
	return &logger
}

func (s *AlarmService) durationUntilNextAction(fallback time.Duration) time.Duration {
	// Check if there is a next alarm action
	logger := log.With().Dur("fallback", fallback).Interface("next_alarm_with_action", s.nextAlarmWithAction).Logger()
	if s.nextAlarmWithAction == nil {
		logger.Debug().Msg("No next alarm action. Returning fallback value.")
		return fallback
	}

	// Calculate time to next action
	now := time.Now()
	until := s.nextAlarmWithAction.NextActionTime.Sub(now)
	if until < 0 {
		logger.Debug().Msg("Next action already passed. Returning fallback value.")
	}
	return until
}

func (s *AlarmService) handleManagerAction(action ManagerAction) *time.Duration {
	log.Debug().Stringer("action", action).Msg("Handling manager action")
	switch action {
	case ManagerActionAbortAlarm:
		s.stopAlarm()
		return nil
	case ManagerActionUpdateSchedule:
		return s.handleUpdateSchedule()
	}
	return nil // Shouldn't be called. Enforced by exhaustive check.
}

func (s *AlarmService) stopAlarm() {
	// Check if an alarm is ringing
	if s.status.IsIdle() {
		log.Warn().Msg("Manager.stopAlarm: Called but no alarm is ringing")
		return
	}

	// Stop alarm and ringer
	alarmID := s.status.GetAlarmID()
	s.status.SetIdle()
	s.ringer.Stop()

	// Disable alarm if not repeated
	var ringingAlarmUpdated bool
	alarm, err := s.GetAlarm(alarmID)
	if err != nil {
		log.Error().Err(err).Stringer("alarm_id", alarmID).Msg("Failed to get alarm")
		return
	}
	if len(alarm.Days) == 0 && alarm.Enabled {
		alarm.Enabled = false
		ringingAlarmUpdated = true
	}
	if err = s.UpdateAlarm(alarm); err != nil {
		log.Error().Err(err).Stringer("alarm_id", alarmID).Msg("Failed to update alarm")
	}

	// Alarm was not updated => Force update of next alarms
	if !ringingAlarmUpdated {
		alarms, err := s.ListAlarms()
		if err != nil {
			log.Error().Err(err).Msg("Failed to list alarms to calculate next alarms")
			return
		}
		nextAlarmToRing, nextAlarmWithAction := s.calculateNextAlarms(alarms)
		s.pubSub.Publish(&pubsub.EventNextAlarmsUpdated{
			ToRing:     nextAlarmToRing,
			WithAction: nextAlarmWithAction,
		})
	}
}

func (s *AlarmService) handleButtonPressed() {
	if s.status.IsRinging() {
		// Alarm is ringing => Stop alarm
		s.stopAlarm()
	} else {
		// No alarm is ringing => Handle night light
		if s.physicalService.GetLEDState().IsOff() {
			s.physicalService.SetLEDState(physicalEntities.LEDState{
				Color:      physicalEntities.PresetColorWarmWhite,
				Brightness: 100,
			})
		} else {
			s.physicalService.UnlockBacklightBrightness()
			s.physicalService.ResetLEDState()
		}
	}
}

func (s *AlarmService) handleButtonLongPressed() {
	if s.status.IsRinging() {
		// Alarm is ringing => Stop alarm
		s.stopAlarm()
	} else {
		// No alarm is ringing => Handle night light
		if s.physicalService.GetLEDState().IsOff() {
			s.physicalService.LockBacklightBrightness()
			s.physicalService.SetLEDState(physicalEntities.LEDState{
				Color:      physicalEntities.PresetColorOrange,
				Brightness: 10,
			})
		} else {
			s.physicalService.UnlockBacklightBrightness()
			s.physicalService.ResetLEDState()
		}
	}
}

func (s *AlarmService) handleUpdateSchedule() *time.Duration {
	// Skip update if not idle
	// Will be updated after alarm is stopped
	if s.status.IsRinging() {
		return nil
	}

	// Skip reschedule if no next alarm
	if s.nextAlarmWithAction == nil {
		return nil
	}

	// Calculate duration
	return lo.ToPtr(s.durationUntilNextAction(time.Second))
}

func (s *AlarmService) handleNextAction() {
	// Check if nextAlarmWithAction is set and is due
	caller := "Manager.handleNextAction"
	if s.nextAlarmWithAction == nil || time.Now().Before(s.nextAlarmWithAction.NextActionTime) {
		// No action required
		log.Warn().Interface("next_alarm_with_action", s.nextAlarmWithAction).Msgf("%s: Called but nextAlarmWithAction not set are not yet due.", caller)
		return
	}

	// Handle next action
	switch s.nextAlarmWithAction.NextAction {
	case entities.NextActionSkip:
		alarm, err := s.GetAlarm(s.nextAlarmWithAction.ID)
		if err != nil {
			log.Error().Stringer("alarm_id", s.nextAlarmWithAction.ID).Msgf("%s: Failed to get alarm to update SkipNext", caller)
			return
		}
		alarm.SkipNext = false
		err = s.UpdateAlarm(alarm)
		if err != nil {
			log.Error().Stringer("alarm_id", s.nextAlarmWithAction.ID).Msgf("%s: Failed to set alarm.SkipNext to false", caller)
			return
		}
	case entities.NextActionRing:
		log.Info().Stringer("alarm_id", s.nextAlarmWithAction.ID).Msgf("%s: Starting alarm", caller)
		s.status.SetRinging(s.nextAlarmWithAction.ID)
		s.lastRings[s.nextAlarmWithAction.ID] = s.nextAlarmWithAction.AlarmTime
		s.ringer.Start()
	}
}
