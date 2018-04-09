package main

func insertSortDesc(sort []int) []int {
	for i := 0; i < len(sort); i++ {
		max := sort[i]
		for j := i + 1; j < len(sort); j++ {
			if sort[j] > max {
				max = sort[j]
				sort[j] = sort[i]
				sort[i] = max
			}
		}
	}
	return sort
}

func insertSortDesc2(sort []int) []int {
	for j := 1; j < len(sort); j++ {
		max := sort[j]
		i := j - 1
		for i > 0 && sort[i] < max {
			sort[i+1] = sort[i]
			i--
		}

		sort[i+1] = max
	}

	return sort
}

// func main() {
// 	sort := []int{1, 2, 3, 4, 5, 6, 7, 8}
// 	fmt.Println(insertSortDesc(sort))
// 	fmt.Println(insertSortDesc2(sort))
// }
