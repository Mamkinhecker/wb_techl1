package main

import (
	"fmt"
	"math/big"
)

func Solve(a, b *big.Int) {
	var sum, sub, mul, div big.Int

	sum.Add(a, b)
	sub.Sub(a, b)
	mul.Mul(a, b)
	div.Div(a, b)

	fmt.Printf("a + b = %v\n", &sum)
	fmt.Printf("a - b = %v\n", &sub)
	fmt.Printf("a * b = %v\n", &mul)
	fmt.Printf("a / b = %v\n", &div)

}
