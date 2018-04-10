package main

import (
	"fmt"
)

func mergeSort(intArray []int, p, r int) {
	if p < r {
		q := (p + r) / 2

		mergeSort(intArray, p, q)
		mergeSort(intArray, q+1, r)
		merge(intArray, p, q, r)
	}
}

func merge(intArray []int, p, q, r int) {

	n1 := q - p + 1
	n2 := r - q

	leftArr := make([]int, n1)
	rightArr := make([]int, n2)
	for i := 0; i < n1; i++ {
		leftArr[i] = intArray[p+i]
	}

	for i := 0; i < n2; i++ {
		rightArr[i] = intArray[q+i+1]
	}
	fmt.Println(leftArr, rightArr)

	i, j := 0, 0

	for k := p; k <= r; k++ {
		if i >= n1 {
			intArray[k] = rightArr[j]
			j++
		} else if j >= n2 {
			intArray[k] = leftArr[i]
			i++
		} else if leftArr[i] <= rightArr[j] {
			intArray[k] = leftArr[i]
			i++
		} else {
			intArray[k] = rightArr[j]
			j++
		}
	}
}

func main() {
	intarray := []int{9, 8, 7, 6, 5}
	mergeSort(intarray, 0, 4)
	fmt.Println(intarray)
	// merge(intarray, 3, 3, 4)
}
