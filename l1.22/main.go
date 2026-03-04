package main

import (
	"fmt"
	"math/big"
)

func Solve(a, b *big.Int) {
	var c big.Int
	c.Add(a, b)
	fmt.Printf("a + b =%v\n", c)

	c.Div(a, b)
	fmt.Printf("a / b =%v\n", c)

	c.Mul(a, b)
	fmt.Printf("a * b =%v\n", c)

	c.Sub(a, b)
	fmt.Printf("a - b =%v\n", c)

}
