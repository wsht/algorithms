package main

import (
	"fmt"
)

/**
LDR
*/
func InOrderTraverse(t *BitTree) {
	if t != nil {
		InOrderTraverse(t.LChild)
		fmt.Println(t.Data)
		InOrderTraverse(t.RChild)
	}
}
