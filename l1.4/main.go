package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	nWorkers := flag.Int("workers", 4, "number of workers")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataChan := make(chan string)
	var wg sync.WaitGroup
	for i := 0; i < *nWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for data := range dataChan {
				fmt.Printf("Worker %d: %s\n", id, data)
			}
			fmt.Printf("Worker %d stopped\n", id)
		}(i)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, shutting down...")
		cancel()
		close(dataChan)
	}()

	go func() {
		counter := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				data := fmt.Sprintf("Data %d", counter)
				dataChan <- data
				counter++
				time.Sleep(time.Second)
			}
		}
	}()

	wg.Wait()
	fmt.Println("All workers stopped, program terminated")
}
