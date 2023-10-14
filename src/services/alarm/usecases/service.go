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
	db              repositories.DB
	pubSub          pubsub.PubSub
	ringer          *Ringer
	physicalService physical.Service

	status             Status                          // Managed by manager
	planningsByAlarmID map[uuid.UUID]entities.Planning // Managed by manager
	timersByAlarmID    map[uuid.UUID]*time.Timer       // Managed by manager
	timerChan          chan uuid.UUID                  // Managed by manager
}

func NewAlarmService(physicalService physical.Service, audioService audio.Service, pubSub pubsub.PubSub) (*AlarmService, error) {
	// Setup DB
	dataPath := "alarms.json"
	db, err := jsondb.NewJSONDB(dataPath)
	if err != nil {
		log.Fatal().Err(err).Str("data_path", dataPath).Msg("Failed to open JSON DB")
	}

	// Build service
	s := &AlarmService{
		db:              db,
		pubSub:          pubSub,
		ringer:          NewRinger(physicalService, audioService),
		physicalService: physicalService,
	}
	return s, s.startManager()
}

func (s *AlarmService) Close() error {
	return s.db.Close()
}
