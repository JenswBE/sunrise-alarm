package usecases

import (
	"time"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/JenswBE/sunrise-alarm/src/services/alarm/planner"
	physicalEntities "github.com/JenswBE/sunrise-alarm/src/services/physical/entities"
	"github.com/JenswBE/sunrise-alarm/src/utils/pubsub"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

func (s *AlarmService) startManager() error {
	// Calculate initial plannings
	alarms, err := s.ListAlarms()
	if err != nil {
		return err
	}
	enabledAlarms := lo.Filter(alarms, func(a entities.Alarm, _ int) bool { return a.Enabled })
	s.planningsByAlarmID = lo.SliceToMap(enabledAlarms, func(a entities.Alarm) (uuid.UUID, entities.Planning) {
		return a.ID, planner.CalculatePlanning(a, time.Now())
	})
	log.Debug().Time("next_ring_time", s.GetNextRingTime()).Msg("Manager.startManager: Initial plannings calculated")

	// Create timers
	s.timerChan = make(chan uuid.UUID)
	s.timersByAlarmID = make(map[uuid.UUID]*time.Timer, len(s.planningsByAlarmID))
	for alarmID, planning := range s.planningsByAlarmID {
		s.timersByAlarmID[alarmID] = createTimer(alarmID, planning, s.timerChan)
	}

	// Start manager loop
	go s.eventLoop()

	return nil
}

func createTimer(alarmID uuid.UUID, planning entities.Planning, timerChan chan uuid.UUID) *time.Timer {
	// Calculate next time
	nextTime := planning.NextRingTime
	if !planning.NextSkipTime.IsZero() {
		nextTime = planning.NextSkipTime
	}

	// Create timer
	return time.AfterFunc(time.Until(nextTime), func() { timerChan <- alarmID })
}

// eventLoop handles all manager events.
// Having a function in a single go-routine handle all data manipulations,
// we are sure we don't create race conditions and using locks is not needed.
func (s *AlarmService) eventLoop() {
	events := make(chan pubsub.Event, 1)
	s.pubSub.Subscribe(pubsub.EventAlarmChanged{}, events)
	s.pubSub.Subscribe(pubsub.EventButtonPressedShort{}, events)
	s.pubSub.Subscribe(pubsub.EventButtonPressedLong{}, events)
	for {
		select {
		case event := <-events:
			switch e := event.(type) {
			case pubsub.EventAlarmChanged:
				log.Debug().Stringer("action", e.Action).Stringer("alarm_id", e.Alarm.ID).Msg("Manager.eventLoop: EventAlarmChanged event received, update plannings and timers...")
				s.handleAlarmChanged(e)
			case pubsub.EventButtonPressedShort:
				log.Debug().Msg("Manager.eventLoop: EventButtonPressedShort event received, handling button press...")
				s.handleButtonPressed()
			case pubsub.EventButtonPressedLong:
				log.Debug().Msg("Manager.eventLoop: EventButtonPressedLong event received, handling button press...")
				s.handleButtonLongPressed()
			}
		case alarmID := <-s.timerChan:
			// Delay will be set by handle_action (through UpdateSchedule)
			log.Debug().Msg("Manager.eventLoop: Received trigger from timer, calling handleTimer...")
			s.handleTimer(alarmID)
		}
	}
}

func (s *AlarmService) handleAlarmChanged(event pubsub.EventAlarmChanged) {
	alarmID := event.Alarm.ID
	logger := log.With().Stringer("action", event.Action).Stringer("alarm_id", alarmID).Logger()
	switch event.Action {
	case pubsub.AlarmChangedActionCreated:
		// Check alarm is enabled
		if !event.Alarm.Enabled {
			logger.Info().Msg("Manager.handleAlarmChanged: Created alarm is not enabled, ignoring event.")
			return
		}

		// Create planning and timer
		logger.Debug().Msg("Manager.handleAlarmChanged: New alarm created, creating planning and timer...")
		planning := planner.CalculatePlanning(event.Alarm, time.Now())
		s.planningsByAlarmID[alarmID] = planning
		s.timersByAlarmID[alarmID] = createTimer(alarmID, planning, s.timerChan)
		logger.Debug().Object("planning", planning).Msg("Manager.handleAlarmChanged: Timer created for planning")
	case pubsub.AlarmChangedActionUpdated:
		// Check alarm is enabled
		if !event.Alarm.Enabled {
			// Might not exist if alarm was disabled and something else changed
			logger.Debug().Msg("Manager.handleAlarmChanged: Updated alarm is disabled, deleting planning and timer (might not exist).")
			s.deletePlanningAndTimer(alarmID, false)
			return
		}

		// Recreate planning and alarm
		logger.Debug().Msg("Manager.handleAlarmChanged: Updated alarm is enabled, recreating planning and timer...")
		s.deletePlanningAndTimer(alarmID, true)
		planning := planner.CalculatePlanning(event.Alarm, time.Now())
		s.planningsByAlarmID[alarmID] = planning
		s.timersByAlarmID[alarmID] = createTimer(alarmID, planning, s.timerChan)
		logger.Debug().Object("planning", planning).Msg("Manager.handleAlarmChanged: Timer created for planning")
	case pubsub.AlarmChangedActionDeleted:
		logger.Debug().Msg("Manager.handleAlarmChanged: Alarm deleted, deleting planning and timer...")
		s.deletePlanningAndTimer(alarmID, true)
	}
}

func (s *AlarmService) deletePlanningAndTimer(alarmID uuid.UUID, expectedToExist bool) {
	// Delete planning
	delete(s.planningsByAlarmID, alarmID)

	// Get planning and timer
	timer, ok := s.timersByAlarmID[alarmID]
	if !ok {
		if expectedToExist {
			log.Error().Stringer("alarm_id", alarmID).Msg("Manager.deletePlanningAndTimer: No timer found for deleted alarm, ignoring timer.")
		}
		return
	}

	// Delete timer
	timer.Stop()
	delete(s.timersByAlarmID, alarmID)
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
	alarm, err := s.GetAlarm(alarmID)
	if err != nil {
		log.Error().Err(err).Stringer("alarm_id", alarmID).Msg("Manager.stopAlarm: Failed to get alarm")
		return
	}
	if len(alarm.Days) == 0 && alarm.Enabled {
		alarm.Enabled = false
	}
	if err = s.UpdateAlarm(alarm); err != nil {
		log.Error().Err(err).Stringer("alarm_id", alarmID).Msg("Manager.stopAlarm: Failed to update alarm")
	}
}

func (s *AlarmService) handleTimer(alarmID uuid.UUID) {
	// Get planning
	logger := log.With().Stringer("alarm_id", alarmID).Logger()
	planning, ok := s.planningsByAlarmID[alarmID]
	if !ok {
		logger.Error().Msg("Manager.handleTimer: Failed to get planning. Ignoring timer.")
		return
	}

	// Handle next action
	switch {
	case !planning.NextSkipTime.IsZero(): // => Skipping alarm
		alarm, err := s.GetAlarm(alarmID)
		if err != nil {
			logger.Error().Err(err).Msg("Manager.handleTimer: Failed to get alarm to update SkipNext")
			return
		}
		alarm.SkipNext = false
		err = s.UpdateAlarm(alarm)
		if err != nil {
			logger.Error().Err(err).Msg("Manager.handleTimer: Failed to set alarm.SkipNext to false")
			return
		}
	default: // => Ringing alarm
		logger.Info().Msg("Manager.handleTimer: Starting alarm")
		s.status.SetRinging(alarmID)
		s.ringer.Start()
	}
}
