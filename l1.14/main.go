package main

import "fmt"

func main() {
	var v interface{}

	fmt.Scan(&v)

	switch i := v.(type) {
	case int:
		fmt.Print(i)
	case string:
		fmt.Print(i)
	case bool:
		fmt.Print(i)
	case chan string:
		fmt.Print(i)
	}

}
