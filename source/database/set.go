package database

import (
	"sort"
	"sync"
)

type (
	// SortedSetEntry represents a singel sorted set entry.
	SortedSetEntry struct {
		member string
		score  float64
		rank   int64
	}

	// SortedSet represents a collection of sorted set entries.
	SortedSet struct {
		sort.Interface
		mutex    sync.RWMutex
		sorted   []*SortedSetEntry
		entries  map[string]*SortedSetEntry
		isSorted bool
	}
)

// CreateSortedSet creates a new, empty, sorted set.
func CreateSortedSet() *SortedSet {
	return &SortedSet{
		sorted:  make([]*SortedSetEntry, 0),
		entries: make(map[string]*SortedSetEntry),
	}
}

// CreateSortedSetEntry creates a single sorted set entry object.
func CreateSortedSetEntry(member string, score float64) (entry *SortedSetEntry) {
	return &SortedSetEntry{
		member: member,
		score:  score,
	}
}

// Add adds a new entry to the sorted set.
func (set *SortedSet) Add(member string, score float64) (exists bool) {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	if _, exists = set.entries[member]; exists {
		if set.entries[member].score == score {
			return
		}

		set.entries[member].score = score
	} else {
		var newEntry = &SortedSetEntry{
			member: member,
			score:  score,
		}

		set.entries[member] = newEntry
		set.sorted = append(set.sorted, newEntry)
	}

	set.isSorted = false
	return
}

// AddEntry adds a new entry (cloned from a source entry) to the sorted set.
func (set *SortedSet) AddEntry(entry *SortedSetEntry) (added bool) {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	if _, exists := set.entries[entry.member]; exists {
		if set.entries[entry.member].score == entry.score {
			return
		}

		set.entries[entry.member].score = entry.score
	} else {
		var newEntry = &SortedSetEntry{
			member: entry.member,
			score:  entry.score,
		}

		set.entries[entry.member] = newEntry
		set.sorted = append(set.sorted, newEntry)
		added = true
	}

	set.isSorted = false
	return
}

// Len returns the number of entries in the sorted set.
func (set *SortedSet) Len() int {
	return len(set.sorted)
}

// Less is used to sort the set using Go standard library sort functions.
func (set *SortedSet) Less(index1, index2 int) bool {
	if set.sorted[index1].score == set.sorted[index2].score {
		return set.sorted[index1].member < set.sorted[index2].member
	}

	return set.sorted[index1].score < set.sorted[index2].score
}

// Swap is used to sort the set using Go standard library sort functions.
func (set *SortedSet) Swap(index1, index2 int) {
	set.sorted[index1], set.sorted[index2] = set.sorted[index2], set.sorted[index1]
}

func (set *SortedSet) checkSort() {
	if !set.isSorted {
		sort.Sort(set)

		for index := range set.sorted {
			set.sorted[index].rank = int64(index)
		}

		set.isSorted = true
	}
}

// Get returns a entry from the sorted set.
func (set *SortedSet) Get(index int) *SortedSetEntry {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	if index < len(set.sorted) {
		set.checkSort()
		return set.sorted[index]
	}

	return nil
}

// GetRank returns the rank for a sorted set member.
func (set *SortedSet) GetRank(member string) int64 {
	if entry, exists := set.entries[member]; exists {
		set.checkSort()
		return entry.rank
	}

	return -1
}

// Get returns the member and score for a sorted set entry.
func (entry *SortedSetEntry) Get() (member string, score float64) {
	return entry.member, entry.score
}

// GetMember returns the member for a sorted set entry.
func (entry *SortedSetEntry) GetMember() string {
	return entry.member
}

// GetScore returns the score for a sorted set entry.
func (entry *SortedSetEntry) GetScore() float64 {
	return entry.score
}
