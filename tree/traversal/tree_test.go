package main

import (
	"fmt"
	"testing"
)

func main() {
	t := &BitTree{}
	t = CreateTree(t, 1, 7)

	fmt.Println(t)
}

func TestCreateTree(x *testing.T) {
	t := &BitTree{}
	t = CreateTree(t, 1, 7)
	fmt.Println(*t)
	// x.Log(t)
	fmt.Println("pre order traverse")
	PreOrderTraverse(t)

	fmt.Println("in order Traverse")
	InOrderTraverse(t)

	fmt.Println("poster order traverse")
	PostOrderTraverse(t)
}
