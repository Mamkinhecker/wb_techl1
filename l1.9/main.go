package main

import (
	"fmt"
	"math/rand"
)

func main() {
	nums := make([]int, 10)
	c1 := make(chan int)
	c2 := make(chan int)

	for i := range 10 {
		nums[i] = rand.Int()
	}

	go func() {
		for _, num := range nums {
			c1 <- num
		}
	}()

	go func() {
		for num := range c1 {
			c2 <- num * num
		}
	}()

	for x := range c2 {
		fmt.Println(x)
	}

}
