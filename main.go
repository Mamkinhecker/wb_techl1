package main

import "fmt"

func main() {
	var s []string
	var n int
	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		var str string
		fmt.Scan(&str)
		s = append(s, str)
	}

	mp := make(map[string]int)

	for _, word := range s {
		mp[word]++
	}

	for word := range mp {
		fmt.Print(word + " ")
	}
}
