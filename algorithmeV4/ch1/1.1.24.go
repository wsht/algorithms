package main

import (
	"fmt"
)

func gcd(p, q int) int {
	fmt.Println(p, q)
	if q == 0 {
		return p
	}

	r := p % q
	return gcd(q, r)
}

// func main() {
// 	fmt.Println(gcd(1111111, 1234567))
// }
