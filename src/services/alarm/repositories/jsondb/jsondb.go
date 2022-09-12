package jsondb

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/JenswBE/sunrise-alarm/src/services/alarm/repositories"
	"github.com/JenswBE/sunrise-alarm/src/services/alarm/repositories/jsondb/internal"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

var _ repositories.DB = &JSONDB{}

type JSONDB struct {
	mutex  sync.RWMutex
	file   *os.File
	alarms map[uuid.UUID]internal.Alarm
}

func NewJSONDB(filePath string) (*JSONDB, error) {
	// Create DB
	db := &JSONDB{alarms: map[uuid.UUID]internal.Alarm{}}

	// Open DB file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON DB file: %w", err)
	}
	db.file = file

	// Get file stats
	fileStats, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get JSON DB file stats: %w", err)
	}

	// Load alarms
	if fileStats.Size() > 0 {
		db.mutex.Lock()
		defer db.mutex.Unlock()
		err = json.NewDecoder(db.file).Decode(&db.alarms)
		if err != nil {
			return nil, fmt.Errorf("failed to load alarms from JSON DB file: %w", err)
		}
	} else {
		// File was just created and therefore empty
		db.alarms = map[uuid.UUID]internal.Alarm{}
		err = db.writeToFile()
		if err != nil {
			return nil, fmt.Errorf("failed to write initial content to JSON DB file: %w", err)
		}
	}

	// Create DB successful
	return db, nil
}

// writeToFile saves the alarms to disk
func (db *JSONDB) writeToFile() error {
	_, err := db.file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to jump to start of JSON DB file: %w", err)
	}
	err = db.file.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed to truncate JSON DB file: %w", err)
	}
	err = json.NewEncoder(db.file).Encode(db.alarms)
	if err != nil {
		return fmt.Errorf("failed to encode alarms and write to JSON DB file: %w", err)
	}
	err = db.file.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync JSON DB file to disk: %w", err)
	}
	return nil
}

func (db *JSONDB) Close() error {
	return db.file.Close()
}

func (db *JSONDB) List() ([]entities.Alarm, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	alarmToGlobalFunc := func(_ uuid.UUID, a internal.Alarm) entities.Alarm { return a.ToGlobal() }
	return lo.MapToSlice(db.alarms, alarmToGlobalFunc), nil
}

func (db *JSONDB) Get(id uuid.UUID) (entities.Alarm, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	alarm, err := db.getWithoutLock(id)
	if err != nil {
		return entities.Alarm{}, err
	}
	return alarm.ToGlobal(), nil
}

func (db *JSONDB) getWithoutLock(id uuid.UUID) (internal.Alarm, error) {
	alarm, found := db.alarms[id]
	if !found {
		return internal.Alarm{}, fmt.Errorf("alarm not found: %s", id.String())
	}
	return alarm, nil
}

func (db *JSONDB) Create(alarm entities.Alarm) (entities.Alarm, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	alarm.ID = uuid.New()
	internalAlarm := internal.AlarmFromGlobal(alarm)
	db.alarms[alarm.ID] = internalAlarm
	return alarm, db.writeToFile()
}

func (db *JSONDB) Update(alarm entities.Alarm) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check alarm exists
	internalAlarm := internal.AlarmFromGlobal(alarm)
	_, err := db.getWithoutLock(internalAlarm.ID)
	if err != nil {
		return err // Very likely alarm doesn't exist
	}

	// Update alarm
	db.alarms[alarm.ID] = internalAlarm
	return db.writeToFile()
}

func (db *JSONDB) Delete(id uuid.UUID) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check alarm exists
	_, err := db.getWithoutLock(id)
	if err != nil {
		return err // Very likely alarm doesn't exist
	}

	// Delete alarm
	delete(db.alarms, id)
	return db.writeToFile()
}
