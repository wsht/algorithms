package main

import (
	"fmt"
	"testing"
)

func TestMergeSort(t *testing.T) {
	arr := []int{
		55, 94, 87, 1, 4, 32, 11, 77, 39, 42, 64, 53, 70, 12, 9,
	}

	fmt.Println("merge sort before arr list:", arr)
	fmt.Println("merge sort before arr list len is:", len(arr))
	MergeSort(arr, 0, len(arr)-1)
	fmt.Println("merge sort after arr list:", arr)
}
