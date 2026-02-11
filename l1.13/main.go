package main

import "fmt"

func main() {
	a, b := 3, 5
	a += b
	b = a - b
	a -= b
	fmt.Print(&a, " ", &b)

}
