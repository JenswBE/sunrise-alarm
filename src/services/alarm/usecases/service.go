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
)

type AlarmService struct {
	db                 repositories.DB
	pubSub             pubsub.PubSub
	ringer             *Ringer
	alarmLightDuration time.Duration
	physicalService    physical.Service

	status              Status                  // Managed by manager
	lastRings           map[uuid.UUID]time.Time // Managed by manager
	nextAlarmToRing     *entities.NextAlarm     // Managed by manager
	nextAlarmWithAction *entities.NextAlarm     // Managed by manager
}

func NewAlarmService(physicalService physical.Service, audioService audio.Service, pubSub pubsub.PubSub, alarmLightDuration, alarmSoundDuration time.Duration) *AlarmService {
	// Setup DB
	dataPath := "alarms.json"
	db, err := jsondb.NewJSONDB(dataPath)
	if err != nil {
		log.Fatal().Err(err).Str("data_path", dataPath).Msg("Failed to open JSON DB")
	}

	// Build service
	managerActions := make(chan ManagerAction, 1)
	s := &AlarmService{
		db:                 db,
		pubSub:             pubSub,
		ringer:             NewRinger(physicalService, audioService, alarmLightDuration, alarmSoundDuration, managerActions),
		alarmLightDuration: alarmLightDuration,
		physicalService:    physicalService,

		lastRings: map[uuid.UUID]time.Time{},
	}
	s.startManager(managerActions)
	return s
}

func (s *AlarmService) Close() error {
	return s.db.Close()
}
