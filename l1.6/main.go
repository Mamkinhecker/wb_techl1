package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup
	stop := false

	wg.Add(1)
	go func() {
		defer wg.Done()
		for !stop {
			fmt.Println("Горутина работает...")
			time.Sleep(500 * time.Millisecond)
		}
		fmt.Println("Горутина остановлена по флагу")
	}()

	time.Sleep(2 * time.Second)
	stop = true
	wg.Wait()

	done := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				fmt.Println("Горутина остановлена по сигналу из канала")
				return
			default:
				fmt.Println("Горутина работает...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	time.Sleep(2 * time.Second)
	close(done)
	wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Горутина остановлена по сигналу из канала")
				return
			default:
				fmt.Println("Горутина работает...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	time.Sleep(2 * time.Second)
	cancel()
	wg.Wait()

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Горутина остановлена по сигналу из канала")
				return
			default:
				fmt.Println("Горутина работает...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Вложенная горутина
		go func() {
			time.Sleep(1 * time.Second)
			fmt.Println("Завершаем горутину изнутри...")
			runtime.Goexit()
		}()

		// Основная горутина продолжает работать
		for i := 0; i < 5; i++ {
			fmt.Println("Основная горутина:", i)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	wg.Wait()

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Горутина восстановлена после panic:", r)
			}
			wg.Done()
		}()

		for i := 0; i < 5; i++ {
			if i == 3 {
				panic("критическая ошибка!")
			}
			fmt.Println("Горутина работает...", i)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	wg.Wait()

	timeout := time.After(2 * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-timeout:
				fmt.Println("Время вышло!")
				return
			default:
				fmt.Println("Горутина работает...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	wg.Wait()

	ticker := time.NewTicker(500 * time.Millisecond)
	done1 := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done1:
				ticker.Stop()
				fmt.Println("Тикер остановлен")
				return
			case t := <-ticker.C:
				fmt.Println("Тик в", t)
			}
		}
	}()

	time.Sleep(2 * time.Second)
	done1 <- true
	wg.Wait()
}
