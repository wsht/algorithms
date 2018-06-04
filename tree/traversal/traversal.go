package main

import "fmt"

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
