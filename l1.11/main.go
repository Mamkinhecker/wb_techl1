package main

import (
	"fmt"
	"slices"
)

func main() {
	a := []int{1, 2, 3}
	b := []int{2, 3, 4}

	for _, i := range a {
		if slices.Contains(b, i) {
			fmt.Print(i)
		}
	}
}
