package main

import "fmt"

/**
DLR
前序遍历
先访问根节点， 再前序访问左节点，最后再前序访问右节点  根-左-右
*/

func PreOrderTraverse(t *BitTree) {
	if t != nil {
		println(t.Data)
		PreOrderTraverse(t.LChild)
		PreOrderTraverse(t.RChild)
	}
}

func PreOrderTraverseWithoutRecursion(t *BitTree) {
	stack := []*BitTree{}
	tmp := t
	//将左子树压入栈后，出栈，再获取右子树，然后再压入右子树的左子树
	for tmp != nil || len(stack) > 0 {
		for tmp != nil {
			stack = append(stack, tmp)
			fmt.Println(tmp.Data)
			tmp = tmp.LChild
		}

		if len(stack) != 0 {
			tmp = stack[len(stack)-1:][0]
			stack = stack[:len(stack)-1]
			tmp = tmp.RChild
		}
	}
}

/**
LDR
中序访问
先中序访问左子树，然后访问根节点，再中序访问右子树， 即左-中-右
*/
func InOrderTraverse(t *BitTree) {
	if t != nil {
		InOrderTraverse(t.LChild)
		fmt.Println(t.Data)
		InOrderTraverse(t.RChild)
	}
}

func InorderTraverseWithoutRecursion(t *BitTree) {
	stack := []*BitTree{}
	tmp := t

	if tmp == nil {
		fmt.Println("the tree is null ")

	} else {
		for tmp != nil || len(stack) != 0 {
			for tmp != nil {
				stack = append(stack, tmp)
				tmp = tmp.LChild
			}
			// fmt.Println("stack len", len(stack))
			if len(stack) != 0 {
				tmp = stack[len(stack)-1:][0]
				// fmt.Println(tmp)
				fmt.Println(tmp.Data)
				stack = stack[:len(stack)-1]
				tmp = tmp.RChild
			}
		}
	}
}

/**
LRD
后序访问
先后序访问左子树，然后后序访问右子树，最后访问根节点即左-右-根
*/
func PostOrderTraverse(t *BitTree) {
	if t != nil {
		PostOrderTraverse(t.LChild)
		PostOrderTraverse(t.RChild)
		fmt.Println(t.Data)
	}
}
