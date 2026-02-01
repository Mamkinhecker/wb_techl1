package main

import (
	"math/rand"
	"sync"
	"time"
)

func main() {

	m := make(map[int]int)
	var mu sync.Mutex
	var wg sync.WaitGroup
	timeout := time.After(3 * time.Second)

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			select {
			case <-timeout:
				return
			default:
				mu.Lock()
				key := rand.Intn(100)
				exp := rand.Intn(100)
				m[key] = exp
				mu.Unlock()

				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			}

		}()
	}

	wg.Wait()
}
