package main

import (
	"fmt"
	"sync"
	"time"
)

func Sleep(d time.Duration) {
	start := time.Now()
	for time.Since(start) < d {
	}
}

func Sleep2(d time.Duration) {
	select {
	case <-time.After(d):
	}
}

func main() {
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func(sec int) {
			defer wg.Done()
			Sleep(time.Second * time.Duration(sec))
			fmt.Printf("горутина %d подождала %d секунд\n", sec, sec)
		}(i)
	}
	wg.Wait()
}
