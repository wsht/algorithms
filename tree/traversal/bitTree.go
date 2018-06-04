package main

import (
	"fmt"
)

type BitTree struct {
	Data   int
	LChild *BitTree
	RChild *BitTree
}

func CreateTree(t *BitTree, startValue, deep int) *BitTree {

	fmt.Println(startValue)
	if startValue > deep {
		t = nil
	} else {
		t = &BitTree{}
		t.Data = startValue
		t.LChild = CreateTree(t.LChild, 2*startValue, deep)
		t.RChild = CreateTree(t.RChild, 2*startValue+1, deep)
	}
	return t
}
