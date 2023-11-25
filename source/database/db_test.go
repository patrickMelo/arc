package database

import (
	"sync"
	"testing"
	"time"
)

func TestDatabaseRace(test *testing.T) {
	var testDB = Create()
	var testWait sync.WaitGroup

	testDB.SetSingleValue("raceTest", "0", 0)

	for index := 0; index < 1000000; index++ {
		testWait.Add(1)

		go func() {
			defer testWait.Done()
			testDB.IncrementSingleValue("raceTest")
		}()
	}

	testWait.Wait()

	if testDB.GetSingleValue("raceTest") != "1000000" {
		test.Fail()
	}
}

func TestExpireTime(test *testing.T) {
	var testDB = Create()

	testDB.SetSingleValue("expireTest", "0", time.Now().Unix()+5)

	time.Sleep(time.Second * 6)

	if testDB.Has("expireTest") {
		test.Fail()
	}
}

func TestSortedSets(test *testing.T) {
	var testSet = CreateSortedSet()

	testSet.Add("one", 1)
	testSet.Add("three", 3)
	testSet.Add("five", 5)
	testSet.Add("four", 4)
	testSet.Add("two", 2)

	if (testSet.Get(0).member != "one") || (testSet.Get(0).score != 1) ||
		(testSet.Get(1).member != "two") || (testSet.Get(1).score != 2) ||
		(testSet.Get(2).member != "three") || (testSet.Get(2).score != 3) ||
		(testSet.Get(3).member != "four") || (testSet.Get(3).score != 4) ||
		(testSet.Get(4).member != "five") || (testSet.Get(4).score != 5) {
		test.Fail()
	}
}
