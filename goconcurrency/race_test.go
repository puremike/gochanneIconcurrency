package main_test

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestDataRaceCondition(t *testing.T) {
	// var mu sync.RWMutex
	var wg sync.WaitGroup

	var state int32

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// mu.Lock()
			// state += int32(i)
			// mu.Unlock()

			// using atomic vaue
			atomic.AddInt32(&state, int32(i))
		}(i)
	}

	wg.Wait()
	t.Log("Final state:", state) // ✅ Ensures all writes are done before the test exits
}
