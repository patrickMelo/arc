package database

import (
	"strconv"
	"sync"
	"time"
)

type (
	// Database defines a database object.
	Database struct {
		mutex sync.RWMutex
		data  map[string]*Value
	}
)

// Create creates a new database object.
func Create() *Database {
	return &Database{
		data: make(map[string]*Value),
	}
}

// Get returns a value from the database.
func (db *Database) Get(key string) *Value {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if value, exists := db.data[key]; exists {
		value.mutex.RLock()
		defer value.mutex.RUnlock()

		if (value.expireTime > time.Now().Unix()) || (value.expireTime == 0) {
			return value
		}
	}

	return nil
}

// GetSingleValue returns a single value from the database.
func (db *Database) GetSingleValue(key string) string {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if value, exists := db.data[key]; exists {
		value.mutex.RLock()
		defer value.mutex.RUnlock()

		if (value.dataType == SingleValue) && (value.expireTime > time.Now().Unix()) || (value.expireTime == 0) {
			return value.data.(string)
		}
	}

	return ""
}

// GetSortedSet returns a sorted set value from the database.
func (db *Database) GetSortedSet(key string) *SortedSet {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if value, exists := db.data[key]; exists {
		value.mutex.RLock()
		defer value.mutex.RUnlock()

		if (value.dataType == SortedSetValue) && (value.expireTime > time.Now().Unix()) || (value.expireTime == 0) {
			return value.data.(*SortedSet)
		}
	}

	return nil
}

// Set sets a database value.
func (db *Database) Set(key string, value *Value) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if value, exists := db.data[key]; exists {
		db.data[key].Set(value.dataType, value.data, value.expireTime)
	} else {
		db.data[key] = value
	}
}

// SetSingleValue sets a database value as a single value.
func (db *Database) SetSingleValue(key string, data string, expires int64) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if value, exists := db.data[key]; exists {
		value.Set(SingleValue, data, expires)
	} else {
		db.data[key] = &Value{
			dataType:   SingleValue,
			expireTime: expires,
			data:       data,
		}
	}
}

// SetSortedSet sets a database value as a sorted set.
func (db *Database) SetSortedSet(key string, set *SortedSet, expires int64) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if value, exists := db.data[key]; exists {
		value.Set(SortedSetValue, set, expires)
	} else {
		db.data[key] = &Value{
			dataType:   SortedSetValue,
			expireTime: expires,
			data:       set,
		}
	}
}

// IncrementSingleValue increments an integer single value.
func (db *Database) IncrementSingleValue(key string) (newValue int64, ok bool) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	var value *Value
	var exists bool

	if value, exists = db.data[key]; exists {
		value.mutex.RLock()

		if value.dataType != SingleValue {
			value.mutex.RUnlock()
			return 0, false
		}

		value.mutex.RUnlock()
	} else {
		value = &Value{
			dataType:   SingleValue,
			expireTime: 0,
			data:       "0",
		}

		db.data[key] = value
	}

	value.mutex.RLock()

	if intValue, err := strconv.ParseInt(value.data.(string), 10, 64); err == nil {
		value.mutex.RUnlock()
		intValue++
		value.Set(SingleValue, strconv.FormatInt(intValue, 10), value.expireTime)
		return intValue, true
	}

	value.mutex.RUnlock()
	return 0, false
}

// Unset removes a value from the database.
func (db *Database) Unset(key string) (had bool) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, had = db.data[key]; had {
		delete(db.data, key)
	}

	return
}

// Has returns if a value exists into the database.
func (db *Database) Has(key string) bool {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if value, exists := db.data[key]; exists && (value != nil) {
		value.mutex.RLock()
		defer value.mutex.RUnlock()

		return (value.expireTime > time.Now().Unix()) || (value.expireTime == 0)
	}

	return false
}

// Size returns the number of database entries.
func (db *Database) Size() (size int) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	// TODO: for performance reasons, this should be a ready to read value, controlling sets and unsets in other db functions

	for _, value := range db.data {
		if value == nil {
			continue
		}

		value.mutex.RLock()

		if (value.expireTime > time.Now().Unix()) || (value.expireTime == 0) {
			size++
		}

		value.mutex.RUnlock()
	}

	return
}
