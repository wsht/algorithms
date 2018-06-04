package main

import (
	"fmt"
)

/**
LRD
*/
func PostOrderTraverse(t *BitTree) {
	if t != nil {
		PostOrderTraverse(t.LChild)
		PostOrderTraverse(t.RChild)
		fmt.Println(t.Data)
	}
}
