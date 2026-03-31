package main

import "sync"

func or(channels ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	var once sync.Once

	for _, ch := range channels {
		go func(c <-chan interface{}) {
			<-c
			once.Do(func() {
				close(out)
			})
		}(ch)
	}

	return out
}
