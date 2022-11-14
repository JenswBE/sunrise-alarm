package usecases

import (
	"time"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/JenswBE/sunrise-alarm/src/services/alarm/repositories"
	"github.com/JenswBE/sunrise-alarm/src/services/alarm/repositories/jsondb"
	"github.com/JenswBE/sunrise-alarm/src/services/audio"
	"github.com/JenswBE/sunrise-alarm/src/services/physical"
	"github.com/JenswBE/sunrise-alarm/src/utils/pubsub"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type AlarmService struct {
	db              repositories.DB
	pubSub          pubsub.PubSub
	ringer          *Ringer
	physicalService physical.Service

	alarmLightDuration time.Duration
	alarmSoundDuration time.Duration
	alarmMaxDuration   time.Duration

	status             Status                          // Managed by manager
	planningsByAlarmID map[uuid.UUID]entities.Planning // Managed by manager
	timersByAlarmID    map[uuid.UUID]*time.Timer       // Managed by manager
	timerChan          chan uuid.UUID                  // Managed by manager
}

func NewAlarmService(physicalService physical.Service, audioService audio.Service, pubSub pubsub.PubSub, alarmLightDuration, alarmSoundDuration time.Duration) (*AlarmService, error) {
	// Setup DB
	dataPath := "alarms.json"
	db, err := jsondb.NewJSONDB(dataPath)
	if err != nil {
		log.Fatal().Err(err).Str("data_path", dataPath).Msg("Failed to open JSON DB")
	}

	// Build service
	abortAlarm := make(chan struct{}, 1)
	s := &AlarmService{
		db:              db,
		pubSub:          pubSub,
		ringer:          NewRinger(physicalService, audioService, alarmLightDuration, alarmSoundDuration, abortAlarm),
		physicalService: physicalService,

		alarmLightDuration: alarmLightDuration,
		alarmSoundDuration: alarmSoundDuration,
		alarmMaxDuration:   lo.Max([]time.Duration{alarmLightDuration, alarmSoundDuration}),
	}
	return s, s.startManager(abortAlarm)
}

func (s *AlarmService) Close() error {
	return s.db.Close()
}
