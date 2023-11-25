package database

import (
	"sync"
)

// Value type constants.
const (
	SingleValue = iota
	SortedSetValue
)

type (
	// Value represents a database value.
	Value struct {
		mutex      sync.RWMutex
		dataType   int
		expireTime int64
		data       interface{}
	}
)

// Get returns the current data for the value.
func (value *Value) Get() interface{} {
	value.mutex.RLock()
	defer value.mutex.RUnlock()

	return value.data
}

// GetType returns the current data type for the value.
func (value *Value) GetType() (dataType int) {
	value.mutex.RLock()
	defer value.mutex.RUnlock()

	return value.dataType
}

// GetInformation returns the current type and expire time for the value.
func (value *Value) GetInformation() (dataType int, expires int64) {
	value.mutex.RLock()
	defer value.mutex.RUnlock()

	return value.dataType, value.expireTime
}

// Set defines new type, data and expire time for the value.
func (value *Value) Set(dataType int, data interface{}, expires int64) {
	value.mutex.Lock()
	defer value.mutex.Unlock()

	value.dataType = dataType
	value.data = data
	value.expireTime = expires
}
