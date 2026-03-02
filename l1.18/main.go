package main

import "sync"

type cnt struct {
	mu  sync.Mutex
	res int
}

func main() {
	p := cnt{
		mu:  sync.Mutex{},
		res: 0,
	}
	var wg sync.WaitGroup
	for range 5 {
		wg.Go(func() {

			p.mu.Lock()
			defer p.mu.Unlock()

			p.res++
		})
	}

	wg.Wait()

}
