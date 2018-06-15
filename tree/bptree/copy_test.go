package bptree

import (
	"fmt"
	"testing"
)

func TestCopy(test *testing.T) {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{6, 7, 8}

	copy(slice2[3:], slice1)
	fmt.Println(slice2)
}
