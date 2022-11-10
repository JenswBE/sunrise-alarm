package usecases

import (
	"sync"

	"github.com/google/uuid"
)

type Status struct {
	mutex   sync.RWMutex
	alarmID uuid.UUID
}

func (s *Status) IsIdle() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.alarmID == uuid.Nil
}

func (s *Status) IsRinging() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.alarmID != uuid.Nil
}

func (s *Status) GetAlarmID() uuid.UUID {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.alarmID
}

func (s *Status) SetIdle() {
	s.mutex.Lock()
	s.alarmID = uuid.Nil
	s.mutex.Unlock()
}

func (s *Status) SetRinging(alarmID uuid.UUID) {
	s.mutex.Lock()
	s.alarmID = alarmID
	s.mutex.Unlock()
}
