package main

/**
堆排序底层实现为一维数组，构建一个逻辑近似完全二叉树，并且满堆积的属性：子节点的键值或者索引总是小于（或者大于）它的父节点

由此可知，构建稳定的堆的顶点一定是最大（最小）值。
然后首末位进行交换，并且从新构建稳定的堆
由此反复，最终得到有序的数组序列

*/

func HeapSort(arr []int) {
	m := len(arr)
	s := m / 2
	//生成最大堆 此时最大值在首位 即arr[0]
	for i := s; i > -1; i-- {
		heap(arr, i, m-1)
	}
	for i := m - 1; i > 0; i-- {
		//将最大值置换到最后一位，并且重新生成稳定堆
		arr[i], arr[0] = arr[0], arr[i]
		heap(arr, 0, i-1)
	}
}

func heap(arr []int, i, end int) {
	//左子树节点
	l := 2*i + 1
	if l > end {
		return
	}

	//n 左右子树最大的节点坐标
	n := l
	//右子树节点
	r := 2*i + 2
	if r <= end && arr[r] > arr[l] {
		n = r
	}
	if arr[i] > arr[n] {
		return
	}
	//交换父节点与最大子节点
	arr[n], arr[i] = arr[i], arr[n]
	heap(arr, n, end)
}

/**
生成最大堆的时候，应用于找到数组中的最大值
*/

/**
1.创建一个堆H[0,...n-1]
2.堆首（最大值）和堆尾进行交换
3.把堆的尺寸缩小1， 并调用shift_down(0), 目的是把新的数组顶端数据调整到相应位置
4.重复步骤2， 直到堆的尺寸为1
*/
