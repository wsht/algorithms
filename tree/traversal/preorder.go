package main

/**
DLR
*/
func PreOrderTraverse(t *BitTree) {
	if t != nil {
		println(t.Data)
		PreOrderTraverse(t.LChild)
		PreOrderTraverse(t.RChild)
	}
}
