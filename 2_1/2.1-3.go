package main

import (
	"errors"
)

func testSearch(search []int, num int) (int, error) {
	for index, value := range search {
		if value == num {
			return index, nil
		}
	}

	return 0, errors.New("error")
}
