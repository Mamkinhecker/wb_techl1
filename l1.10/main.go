package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	arr := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	sort.Float64s(arr)
	var last float64 = 0
	for i := range 8 {
		if math.Floor(arr[i]/10) != last {
			if i == 0 {
				fmt.Print(math.Floor(arr[i]/10)*10, "{", arr[i])
			} else {
				fmt.Print("}, ", math.Floor(arr[i]/10)*10, "{", arr[i])
			}
			last = math.Floor(arr[i] / 10)
		} else {
			fmt.Print(", ", arr[i])
		}
	}
	fmt.Print("}")
}
