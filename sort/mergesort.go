package main

func MergeSort(arr []int, start, end int) {
	if end > start {
		// fmt.Println(start, end)
		mid := (end + start) / 2
		MergeSort(arr, start, mid)
		MergeSort(arr, mid+1, end)
		merge(arr, start, mid, end)
	}
}

/**
合并数组
依照从小到大排序
*/
func merge(arr []int, start, mid, end int) {

	// fmt.Println("merge start, mid, end ;", start, mid, end)
	leftLen := mid - start + 1
	rightLen := end - mid

	leftList := make([]int, leftLen)
	rightList := make([]int, rightLen)

	for i := 0; i < leftLen; i++ {
		leftList[i] = arr[start+i]
	}

	for j := 0; j < rightLen; j++ {
		rightList[j] = arr[mid+j+1]
	}
	// fmt.Println("left list:", leftList, "rightList:", rightList)
	i, j := 0, 0
	for k := start; k <= end; k++ {
		//the left list is << right list
		if i >= leftLen {
			arr[k] = rightList[j]
			j++
		} else if j >= rightLen { //the left list is >> right list
			arr[k] = leftList[i]
			i++
		} else if leftList[i] <= rightList[j] {
			arr[k] = leftList[i]
			i++
		} else {
			arr[k] = rightList[j]
			j++
		}
	}
}
