package badgerdb

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/JenswBE/sunrise-alarm/src/srv-config/common"
	"github.com/JenswBE/sunrise-alarm/src/srv-config/entities"
	"github.com/JenswBE/sunrise-alarm/src/srv-config/repositories"
	"github.com/google/uuid"
)

var _ repositories.DB = &JSONDB{}

type JSONDB struct {
	mutex  sync.RWMutex
	file   *os.File
	alarms map[uuid.UUID]entities.Alarm
}

func NewJSONDB(filePath string) (*JSONDB, error) {
	// Create DB
	db := &JSONDB{alarms: map[uuid.UUID]entities.Alarm{}}

	// Open DB file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON DB file: %w", err)
	}
	db.file = file

	// Load alarms
	db.mutex.Lock()
	defer db.mutex.Unlock()
	err = json.NewDecoder(db.file).Decode(&db.alarms)
	if err != nil {
		return nil, fmt.Errorf("failed to load alarms from JSON file: %w", err)
	}

	// Create DB successful
	return db, nil
}

// writeToFile saves the alarms to disk
func (db *JSONDB) writeToFile() error {
	err := json.NewEncoder(db.file).Encode(db.alarms)
	if err != nil {
		return fmt.Errorf("failed to encode alarms and write to JSON file: %w", err)
	}
	err = db.file.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync JSON file to disk: %w", err)
	}
	return nil
}

func (db *JSONDB) Close() error {
	return db.file.Close()
}

func (db *JSONDB) List() ([]entities.Alarm, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	return common.ToSlice(db.alarms), nil
}

func (db *JSONDB) Get(id uuid.UUID) (entities.Alarm, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	return db.getWithoutLock(id)
}

func (db *JSONDB) getWithoutLock(id uuid.UUID) (entities.Alarm, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	alarm, found := db.alarms[id]
	if !found {
		return entities.Alarm{}, fmt.Errorf("alarm not found: %s", id.String())
	}
	return alarm, nil
}

func (db *JSONDB) Create(alarm entities.Alarm) (entities.Alarm, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	alarm.ID = uuid.New()
	db.alarms[alarm.ID] = alarm
	return alarm, nil
}

func (db *JSONDB) Update(alarm entities.Alarm) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check alarm exists
	_, err := db.getWithoutLock(alarm.ID)
	if err != nil {
		return err // Very likely alarm doesn't exist
	}

	// Update alarm
	db.alarms[alarm.ID] = alarm
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
