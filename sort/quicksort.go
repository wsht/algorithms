package main

import (
	"fmt"
	"time"
)

func quickSort(value []int, left, right int) {
	if left < right {
		// fmt.Println("valeue left and right is ", left, right)
		time.Sleep(time.Millisecond * 100)
		p := partition(value, left, right)
		// fmt.Println("middle index is", p)
		// fmt.Println(value)
		quickSort(value, left, p-1)
		quickSort(value, p, right)
	}
}

func partition(a []int, left, right int) int {
	pivot := a[right]
	i := left - 1
	for j := left; j < right; j++ {
		if a[j] < pivot {
			i++
			a[i], a[j] = a[j], a[i]
		}
	}

	a[i+1], a[right] = a[right], a[i+1]
	return i + 1
}

func myunserstandQuickSort(arr []int, head, tail int) {
	if head < tail {
		pivot := arr[head]
		p := head
		i, j := head, tail
		for i < j {
			for j >= p && arr[j] >= pivot {
				j--
			}
			if j >= p {
				arr[p], arr[j] = arr[j], arr[p]
				p = j
			}

			for i <= p && arr[i] <= pivot {
				i++
			}
			if i <= p {
				arr[p], arr[i] = arr[i], arr[p]
				p = i
			}
		}
		fmt.Println("value head and tail is ", head, tail)
		fmt.Println("the middle index is", p)
		fmt.Println(arr)
		if p-head > 1 {
			myunserstandQuickSort(arr, head, p-1)
		}

		if tail-p > 1 {
			myunserstandQuickSort(arr, p+1, tail)
		}
	}
}

func main() {
	var array = []int{5, 2, 9, 1, 7, 6}
	var array2 = []int{5, 2, 9, 1, 7, 6}
	quickSort(array, 0, len(array)-1)
	fmt.Println("quick sort ", array)
	fmt.Println("before quick sort arr is ", array2)
	myunserstandQuickSort(array2, 0, len(array)-1)
	fmt.Println("myunserstand quick sort", array2)
}
