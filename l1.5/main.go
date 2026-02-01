package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	c := make(chan int)
	var n int
	fmt.Scan(&n)
	timer := time.After(time.Duration(n) * time.Second)

	go func() {
		for {
			n := rand.Intn(100)
			c <- n
		}
	}()

	go func() {
		for {
			a := <-c
			fmt.Println(a)
		}
	}()

	<-timer
	fmt.Printf("Программа завершена после %d секунд\n", n)
}
