package flashstore

import (
	"sync"
)

// FlashStore is an in-memory flash store.
// The zero value is ready to use.
// Do not copy a non-zero FlashStore.
type FlashStore[T any] struct {
	mutex  sync.Mutex
	values []T
}

// Add adds a value to the flash store
func (s *FlashStore[T]) Add(value T) {
	s.mutex.Lock()
	if s.values == nil {
		s.values = []T{value}
	} else {
		s.values = append(s.values, value)
	}
	s.mutex.Unlock()
}

// Get clears the flash store and returns all values
func (s *FlashStore[T]) Get() []T {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.values == nil {
		return []T{}
	}
	values := s.values
	s.values = []T{}
	return values
}
