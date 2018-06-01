package main

import (
	"fmt"
	"math"
)

func binomial(n, k int, p float64, deep int) float64 {
	fmt.Println(deep)
	deep++
	if n == 0 && k == 0 {
		return 1.0
	}

	if n < 0 || k < 0 {
		return 0.0
	}
	return (1.0-p)*binomial(n-1, k, p, deep) + p*binomial(n-1, k-1, p, deep)
}

/**
the question is
p^k (1-p)^N-K
*/
func binomail2(n, k int, p float64) float64 {

	binomailMap := make([][]float64, n+1)
	for i := 0; i <= n; i++ {
		tmp := math.Pow(1.0-p, float64(i))
		binomailMap[i] = append(binomailMap[i], tmp)
	}
	binomailMap[0][0] = 1.0

	for i := 0; i <= n; i++ {
		for j := 1; j <= k; j++ {
			tmp := p*binomailMap[i-1][j-1] + (1.0-p)*binomailMap[i-1][j]
			binomailMap[i] = append(binomailMap[i], tmp)
		}
	}

	return binomailMap[n][k]
}

func main() {
	fmt.Println(binomial(5, 5, 0.25, 0))
	fmt.Println(binomail2(5, 5, 0.25))
}
