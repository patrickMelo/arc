package vm

import (
	"strconv"
	"time"

	"arc/database"
)

// SET key value
func stdSet(db *database.Database, parameters []string) []string {
	if len(parameters) != 2 {
		return invalidParametersResult
	}

	db.SetSingleValue(parameters[0], parameters[1], 0)
	return okResult
}

// SET key value EX seconds
func stdSetEx(db *database.Database, parameters []string) []string {
	if (len(parameters) != 4) || (parameters[2] != "EX") {
		return invalidParametersResult
	}

	if expireSeconds, err := strconv.ParseInt(parameters[3], 10, 64); err == nil {
		db.SetSingleValue(parameters[0], parameters[1], time.Now().Unix()+expireSeconds)
		return okResult
	}

	return invalidParameterValueResult
}

// GET key
func stdGet(db *database.Database, parameters []string) []string {
	if len(parameters) != 1 {
		return invalidParametersResult
	}

	if !db.Has(parameters[0]) {
		return nilResult
	}

	return []string{db.GetSingleValue(parameters[0])}
}

// DEL key [key...]
func stdDel(db *database.Database, parameters []string) []string {
	if len(parameters) < 1 {
		return invalidParametersResult
	}

	var delCounter int64

	for index := range parameters {
		if db.Unset(parameters[index]) {
			delCounter++
		}
	}

	return []string{strconv.FormatInt(delCounter, 10)}
}

// DBSIZE
func stdDbSize(db *database.Database, parameters []string) []string {
	if len(parameters) != 0 {
		return invalidParametersResult
	}

	return []string{strconv.FormatInt(int64(db.Size()), 10)}
}

// INCR key
func stdIncr(db *database.Database, parameters []string) []string {
	if len(parameters) != 1 {
		return invalidParametersResult
	}

	if newValue, ok := db.IncrementSingleValue(parameters[0]); ok {
		return []string{strconv.FormatInt(newValue, 10)}
	}

	return invalidDataTypeResult
}

// ZADD key score member [score member...]
func stdZadd(db *database.Database, parameters []string) []string {
	if (len(parameters) < 3) || (len(parameters)%2 == 0) {
		return invalidParametersResult
	}

	var set *database.SortedSet

	if value := db.Get(parameters[0]); value == nil {
		set = database.CreateSortedSet()
		db.SetSortedSet(parameters[0], set, 0)
	} else if value.GetType() != database.SortedSetValue {
		return invalidDataTypeResult
	} else {
		set = value.Get().(*database.SortedSet)
	}

	// Parse the values to make sure they are valid before adding (so we can mimic a transaction like - all or none - operation).

	var numberOfElements = (len(parameters) - 1) / 2
	var entries = make([]*database.SortedSetEntry, numberOfElements)

	for index := 0; index < numberOfElements; index++ {
		if value, err := strconv.ParseFloat(parameters[(index*2)+1], 64); err == nil {
			entries[index] = database.CreateSortedSetEntry(parameters[(index*2)+2], value)
		} else {
			return invalidParameterValueResult
		}
	}

	var addCounter int64

	for index := range entries {
		if set.AddEntry(entries[index]) {
			addCounter++
		}
	}

	return []string{strconv.FormatInt(addCounter, 10)}
}

// ZCARD key
func stdZcard(db *database.Database, parameters []string) []string {
	if len(parameters) != 1 {
		return invalidParametersResult
	}

	if set := db.GetSortedSet(parameters[0]); set != nil {
		return []string{strconv.FormatInt(int64(set.Len()), 10)}
	}

	return []string{"0"}
}

// ZRANK key member
func stdZrank(db *database.Database, parameters []string) []string {
	if len(parameters) != 2 {
		return invalidParametersResult
	}

	if set := db.GetSortedSet(parameters[0]); set != nil {
		if rank := set.GetRank(parameters[1]); rank >= 0 {
			return []string{strconv.FormatInt(rank, 10)}
		}
	}

	return nilResult
}

// ZRANGE key start stop
func stdZrange(db *database.Database, parameters []string) []string {
	if len(parameters) != 3 {
		return invalidParametersResult
	}

	if set := db.GetSortedSet(parameters[0]); set != nil {
		var start, startError = strconv.ParseInt(parameters[1], 10, 64)
		var stop, stopError = strconv.ParseInt(parameters[2], 10, 64)

		if (startError != nil) || (stopError != nil) {
			return invalidParameterValueResult
		}

		var size = int64(set.Len())

		if stop < 0 {
			stop += size
		} else if stop >= size {
			stop = size - 1
		}

		if start < 0 {
			start = 0
		} else if start > stop {
			return emptyResult
		}

		var result = make([]string, stop-start+1)

		for index := start; index <= stop; index++ {
			result[index-start] = set.Get(int(index)).GetMember()
		}

		return result
	}

	return emptyResult
}
