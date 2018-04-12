package main

import (
	"fmt"
	"math"
)

//编写一个函数计算ln(N!)的值
/**
这里没有进行尾递归优化，当n的值特别大的时候，会出现栈溢出
*/
func myLn(n int) float64 {
	if n == 1 {
		return math.Log(1)
	} else {
		return math.Log(float64(n)) + myLn(n-1)
	}
}

func myLn2(n int, sum []float64) float64 {
	sum[0] += math.Log(float64(n))
	if n == 1 {
		return sum[0]
	}
	//go 没有尾递归优化么？
	return myLn2(n-1, sum)
}

func main() {
	sum := []float64{0}
	fmt.Println(myLn2(10000000, sum))
	fmt.Println(myLn(10000000))

}
