package main

import (
	"fmt"
	"math"
)

func addTwoBinaryNums(nums1, nums2 []int) []int {
	carray := 0
	nums1Len := len(nums1)
	nums2Len := len(nums2)

	len := int(math.Max(float64(nums1Len), float64(nums2Len)))
	nums1Between := int(math.Max(float64(len-nums1Len), float64(0)))
	nums2Between := int(math.Max(float64(len-nums2Len), float64(0)))

	result := make([]int, len+1)

	for i := len - 1; i >= 0; i-- {
		sum1 := 0
		if i-nums1Between >= 0 {
			sum1 = nums1[i-nums1Between]
		}

		sum2 := 0
		if i-nums2Between >= 0 {
			sum2 = nums2[i-nums2Between]
		}
		sum := sum1 + sum2 + carray
		result[i+1] = sum % 2
		carray = sum / 2
	}

	if carray == 1 {
		result[0] = 1
		return result
	} else {
		return result[1:]
	}
}

func main() {
	fmt.Println(addTwoBinaryNums([]int{1, 0, 1}, []int{1}))
}
