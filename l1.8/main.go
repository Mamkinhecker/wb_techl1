package main

import "fmt"

func main() {
	var a int64
	fmt.Scan(&a)

	var ind int
	fmt.Scan(&ind)

	var typ int
	fmt.Scan(&typ)

	if typ == 1 {
		a = a | (1 << ind)
	} else {
		a = a & ^(1 << ind)
	}

	fmt.Print(a)
}
