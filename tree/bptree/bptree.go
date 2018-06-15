package bptree

import (
	"fmt"
	"io"
	"sync"
)

const (
	kx = 32 //TODO benchmark tune this number if using custom key/value type(s)
	kd = 32 //TODO benchmark tune this number if using custom key/value type(s)
)

func init() {
	if kd < 1 {
		panic(fmt.Errorf("kd %d: out of range", kd))
	}

	if kx < 2 {
		panic(fmt.Errorf("kx %d: out of range", kx))
	}
}

type btEnumeratorPool struct {
	sync.Pool
}

func (p *btEnumeratorPool) get(err error, hit bool, i int, key interface{}, enumeratorQ *dataPage, tree *Tree, ver int64) *Enumerator {
	x := p.Get().(*Enumerator) //get from enumerator pool
	x.err, x.hit, x.i, x.key, x.enumeratorQ, x.tree, x.ver =
		err, hit, i, key, enumeratorQ, tree, ver
	return x
}

type btTreePool struct {
	sync.Pool
}

func (p *btTreePool) get(cmp Cmp) *Tree {
	x := p.Get().(*Tree)
	x.cmp = cmp
	return x
}

type (
	Cmp func(a, b interface{}) int
	//索引页面结构
	//索引结构页面总 item 的存储顺序是
	//{{ch,key},{ch,key},{ch,key},{ch, key},{ch, nil}}
	//其中 >= item[i].key and < item[i+1].key 向 item[i].ch下继续查找
	indexPage struct {
		count int
		item  [2*kx + 2]indexItem
	}
	//索引实体结构
	indexItem struct {
		key interface{} // index key
		ch  interface{} //首部指针 只想下一层的 indexPage 或者dataPage
	}

	//数据页面结构
	//其中item 是使用slice存储的，在结构上是连续的
	dataPage struct {
		count int
		item  [2*kd + 1]dataItem
		next  *dataPage //下一个数据页面指针
		prev  *dataPage //上一个数据页面指针
	}

	//数据实体结构 保存key value值
	dataItem struct {
		key   interface{}
		value interface{}
	}

	Enumerator struct {
		err         error
		hit         bool
		i           int
		key         interface{}
		tree        *Tree
		ver         int64
		enumeratorQ *dataPage //? //记录当前数据页
	}

	Tree struct {
		count int
		cmp   Cmp
		first *dataPage
		last  *dataPage
		TreeR interface{}
		ver   int64
	}
)

var (
	btDataPagePoolEnity   = sync.Pool{New: func() interface{} { return &dataPage{} }}
	btEnumeratorPoolEnity = btEnumeratorPool{sync.Pool{New: func() interface{} { return &Enumerator{} }}}
	btTreePoolEnity       = btTreePool{sync.Pool{New: func() interface{} { return &Tree{} }}}
	btIndexPagePoolEnity  = sync.Pool{New: func() interface{} { return &indexPage{} }}
)

var (
	zdatapage   dataPage
	zdataitem   dataItem
	zenumerator Enumerator
	zkey        interface{}
	ztree       Tree
	zindexpage  indexPage
	zindexitem  indexItem
)

//清除该节点 以及所有子元素
func clr(q interface{}) {
	switch x := q.(type) {
	case *indexPage:
		for i := 0; i < x.count; i++ {
			clr(x.item[i].ch)
		}
		*x = zindexpage
		btIndexPagePoolEnity.Put(x)
	case *dataPage:
		*x = zdatapage
		btDataPagePoolEnity.Put(x)
	}
}

//----------------------indexPage
//创建一个新的索引页 ch0 头部指针 只想下一元素
func newIndexPage(ch0 interface{}) *indexPage {
	r := btIndexPagePoolEnity.Get().(*indexPage)
	r.item[0].ch = ch0
	return r
}

//提取 或者说是删除 索引页面item中 第 i个索引
//注意 item 的分配空间是 2*kx + 2 个 indexItem
//原count=4
//当kx=2时 i=1 例如 {count:4,item:{0,1,2,3,zindexitem,zindexitem}}
func (q *indexPage) extract(i int) {
	q.count-- //{count:3,item:{0,1,2,3,zindexitem,zindexitem}}
	if i < q.count {
		copy(q.item[i:], q.item[i+1:q.count+1]) //{count:3, item:{0,2,3,3,zindexitem, zindexitem}}
		q.item[q.count].ch = q.item[q.count+1].ch
		q.item[q.count].key = zkey     //gc
		q.item[q.count+1] = zindexitem //gc //{count:3, item:{0,2,3,zindexitem, zindexitem, zindexitem}}
	}
}

//在i位置插入 indexitem:{key, ch}
// 当 kx=3时候，i=1 {count:5, item:{key0,ch0},{key1,ch1},{key2,ch2},{key3,ch3},{key4,ch4},{ch5,nil}, {nil,nil},{nil,nil}}
func (q *indexPage) insert(i int, key interface{}, ch interface{}) *indexPage {
	count := q.count
	if i < count {
		q.item[count+1].ch = q.item[count].ch //{count:5, item:{key0,ch0},{key1,ch1},{key2,ch2},{key3,ch3},{key4,ch4},{ch5,nil}, {ch5,nil},{nil,nil}}
		copy(q.item[i+2:], q.item[i+1:count]) //{count:5, item:{key0,ch0},{key1,ch1},{key2,ch2},{key2,ch2},{key3,ch3},{key4,ch4}, {ch5,nil},{nil,nil}}
		q.item[i+1].key = q.item[i].key       //{count:5, item:{key0,ch0},{key1,ch1},{key1,ch2},{key2,ch2},{key3,ch3},{key4,ch4}, {ch5,nil},{nil,nil}}
	}

	count++
	q.count = count     //{count:6, item:{key0,ch0},{key1,ch1},{key1,ch2},{key2,ch2},{key3,ch3},{key4,ch4}, {ch5,nil},{nil,nil}}
	q.item[i].key = key ////{count:6, item:{key0,ch0},{insertKey,ch1},{key1,ch2},{key2,ch2},{key3,ch3},{key4,ch4}, {ch5,nil},{nil,nil}}
	q.item[i+1].ch = ch //{count:6, item:{key0,ch0},{insertKey,ch1},{key1,insertCh},{key2,ch2},{key3,ch3},{key4,ch4}, {ch5,nil},{nil,nil}}
	return q
}

/**
获取兄弟元素
*/
func (q *indexPage) siblings(i int) (l, r *dataPage) {
	if i >= 0 {
		l = q.item[i-1].ch.(*dataPage)
	}
	if i < q.count {
		r = q.item[i+1].ch.(*dataPage)
	}
	return
}

//----------------------------datapage
//适用范围 当当前datapage页面已经满了，但是其左兄弟没有满的情况下，将其一部分转入左兄弟中
//左旋操作，
func (left *dataPage) mvL(right *dataPage, count int) {
	copy(left.item[left.count:], right.item[:count])
	copy(right.item[:], right.item[count:right.count])
	left.count += count
	right.count -= count
}

//右旋操作
func (left *dataPage) mvR(right *dataPage, count int) {
	copy(right.item[count:], right.item[:count])
	copy(right.item[:count], left.item[left.count-count:])
	left.count -= count
	right.count += count
}

//---------------------------tree
//TreeNew returns a newly created, empty Tree. The compare function is used
//for key collation
func TreeNew(cmp Cmp) *Tree {
	return btTreePoolEnity.get(cmp)
}

func (t *Tree) Clear() {
	if t.TreeR == nil {
		return
	}
	clr(t.TreeR)
	t.count, t.first, t.last, t.TreeR = 0, nil, nil, nil
	t.ver++
}

func (t *Tree) Close() {
	t.Clear()
	*t = ztree
	btTreePoolEnity.Put(t)
}

func (t *Tree) insert(dataPage *dataPage, i int, key interface{}, value interface{}) *dataPage {
	t.ver++
	count := dataPage.count
	if i < count {
		//这里是不是丢失了一位，没有append操作感觉
		copy(dataPage.item[i+1:], dataPage.item[i:count])
	}

	count++
	dataPage.count = count
	dataPage.item[i].key, dataPage.item[i].value = key, value
	t.count++
	return dataPage
}

func (t *Tree) find(q interface{}, key interface{}) (i int, ok bool) {
	var middleKey interface{}
	low := 0
	switch x := q.(type) {
	case *indexPage:
		hight := x.count - 1
		//二分查找法
		for low <= hight {
			middle := (low + hight) >> 1 //  l+h/2
			middleKey = x.item[middle].key
			switch cmp := t.cmp(key, middleKey); {
			case cmp > 0: //k > mk
				low = middle + 1
			case cmp == 0:
				return middle, true
			default:
				hight = middle - 1
			}
		}
	case *dataPage:
		hight := x.count - 1
		for low <= hight {
			middle := (low + hight) >> 1
			middleKey = x.item[middle].key
			switch cmp := t.cmp(key, middleKey); {
			case cmp > 0: // k > mk
				low = middle + 1
			case cmp == 0:
				return middle, true
			default:
				hight = middle - 1
			}
		}
	}

	return low, false
}

func (t *Tree) extract(q *dataPage, i int) {
	t.ver++
	q.count--
	if i < q.count {
		copy(q.item[i:], q.item[i+1:q.count+1])
	}
	q.item[q.count] = zdataitem
	t.count--
}

func (t *Tree) First() (key interface{}, value interface{}) {
	if q := t.first; q != nil {
		q := &q.item[0]
		key, value = q.key, q.value
	}
	return
}

func (t *Tree) Last() (key interface{}, value interface{}) {
	if q := t.last; q != nil {
		q := &q.item[q.count-1]
		key, value = q.key, q.value
	}
	return
}

func (t *Tree) Len() int {
	return t.count
}

func (t *Tree) Get(key interface{}) (value interface{}, ok bool) {
	q := t.TreeR
	if q == nil {
		return
	}
	//todo how find the true value ?
	for {
		var i int
		if i, ok := t.find(q, key); ok {
			switch x := q.(type) {
			case *indexPage:
				q = x.item[i+1].ch
				continue
			case *dataPage:
				return x.item[i].value, true
			}
		}

		//没有找到的情况
		switch x := q.(type) {
		case *indexPage:
			q = x.item[i].ch
		default:
			return
		}
	}
}

/**
Seek returns an Enumerator positioned on an item such that k >= item's key
ok reports if k == item.key The enumerator's position is possibly after the
last item in the tree
*/
func (t *Tree) Seek(key interface{}) (e *Enumerator, ok bool) {
	q := t.TreeR
	if q == nil {
		e = btEnumeratorPoolEnity.get(nil, false, 0, key, nil, t, t.ver)
		return
	}

	for {
		var i int
		if i, ok = t.find(q, key); ok {
			switch x := q.(type) {
			case *indexPage:
				q = x.item[i+1].ch
				continue
			case *dataPage:
				return btEnumeratorPoolEnity.get(nil, ok, i, key, x, t, t.ver), true
			}
		}

		switch x := q.(type) {
		case *indexPage:
			q = x.item[i].ch
		case *dataPage:
			return btEnumeratorPoolEnity.get(nil, ok, i, key, x, t, t.ver), false
		}
	}
}

/**
SeekFirst returns an enumerator positioned on the first KV pair in the tree,
if any. for empty tree, err == io.EOF is returned and e will be nil
*/
func (t *Tree) SeekFirst() (e *Enumerator, err error) {
	q := t.first
	if q == nil {
		return nil, io.EOF
	}

	return btEnumeratorPoolEnity.get(nil, true, 0, q.item[0].key, q, t, t.ver), nil
}

/**
SeekLast returns an enumerator positionsed on the last KV pair in the three,
if any, For an empty tree, err == io.EOF is returned and e will be nil
*/
func (t *Tree) SeekLast() (e *Enumerator, err error) {
	q := t.last
	if q == nil {
		return nil, io.EOF
	}

	return btEnumeratorPoolEnity.get(nil, true, q.count-1, q.item[q.count-1].key, q, t, t.ver), nil
}

/**
B+树的插入以及更新
插入过程：
根据插入值的大小，逐步向下直到对应的叶子节点。
如果椰子节点关键字个数小于2t，则直接插入或者更新卫星数据。
如果插入之前，子节点已经满了，则分裂该叶子节点成两半，并把中间值提上到父节点的关键字中，如果该操作导致父节点满了的话，则把父节点分裂，如此递归向上。
*/
func (t *Tree) Set(key interface{}, value interface{}) {
	pi := -1
	var page *indexPage
	q := t.TreeR
	if q == nil { //空树的插入
		z := t.insert(btDataPagePoolEnity.Get().(*dataPage), 0, key, value)
		t.TreeR, t.first, t.last = z, z, z
		return
	}

	for {
		i, ok := t.find(page, key)
		//key值存在
		if ok {
			switch x := q.(type) {
			case *indexPage:
				i++
				if x.count > 2*kx {
					x, i = t.splitIndexPage(page, x, pi, i)
				}
				pi = i
				page = x
				q = x.item[i].ch
				continue
			case *dataPage: //new KV pair
				x.item[i].value = value
			}

			return
		}

		switch x := q.(type) {
		case *indexPage:
			if x.count > 2*kx {
				x, i = t.splitIndexPage(page, x, pi, i)
			}
			pi = i
			page = x
			q = x.item[i].ch
		case *dataPage:
			switch {
			case x.count < 2*kd:
				t.insert(x, i, key, value)
			default:
				//溢出处理
				t.overflow(page, x, pi, i, key, value)
			}
			return
		}
	}

}

/**
Put combines Get and Set in a more efficient way where the tree is walked
only once. The upd(after) receives(old-value, true) if a KV pair for key
exists or (zero-value, false) otherwise.
It can then return a (new-value, true) to create or overwirte the existing value in the KV pair, or
(whatever, false) if it decides not to create or not to update the value of the KV pair

	tree.Set(key, value) call conceptually equals calling
	tree.Put(key, func(interface{}, bool)){return v, true}

modulo the differing return values.

*/
func (t *Tree) Put(key interface{}, upd func(oldValue interface{}, exists bool) (newValue interface{}, write bool)) (oldValue interface{}, written bool) {
	pi := -1
	var p *indexPage
	q := t.TreeR
	var newValue interface{}
	if q == nil {
		newValue, written = upd(newValue, false)
		if !written {
			return
		}

		z := t.insert(btDataPagePoolEnity.Get().(*dataPage), 0, key, newValue)
		t.TreeR, t.first, t.last = z, z, z
		return
	}

	for {
		i, ok := t.find(q, key)
		if ok {
			switch x := q.(type) {
			case *indexPage:
				i++
				if x.count > 2*kx {
					x, i = t.splitIndexPage(p, x, pi, i)
				}
				pi = i
				p = x
				q = x.item[i].ch
				continue
			case *dataPage:
				oldValue = x.item[i].value
				newValue, written = upd(oldValue, true)
				if !written {
					return
				}
				x.item[i].value = newValue
			}
			return
		}

		switch x := q.(type) {
		case *indexPage:
			if x.count > 2*kx {
				x, i = t.splitIndexPage(p, x, pi, i)
			}
			pi = i
			p = x
			q = x.item[i].ch
		case *dataPage:
			newValue, written = upd(newValue, false)
			if !written {
				return
			}

			switch {
			case x.count < 2*kd:
				t.insert(x, i, key, newValue)
			default:
				t.overflow(p, x, pi, i, key, newValue)
			}
			return
		}
	}
}

/**
Delete removes the key's KV pair, if it exists, in which case Delete returns true
*/
func (t *Tree) Delete(key interface{}) bool {
	pi := -1
	var indexpage *indexPage
	q := t.TreeR
	if q == nil {
		return false
	}

	for {
		var i int
		i, ok := t.find(q, key)
		if ok {
			switch x := q.(type) {
			case *indexPage:
				if x.count < kx && q != t.TreeR {
					x, i = t.underflowIndexPage(indexpage, x, pi, i)
				}
				pi = i + 1
				indexpage = x
				q = x.item[pi].ch
				continue
			case *dataPage:
				t.extract(x, i)
				if x.count >= kd {
					return true
				}

				if q != t.TreeR {
					t.underflow(indexpage, x, pi)
				} else if t.count == 0 {
					t.Clear()
				}

				return true
			}
		}

		switch x := q.(type) {
		case *indexPage:
			if x.count < kx && q != t.TreeR {
				x, i = t.underflowIndexPage(indexpage, x, pi, i)
			}

			pi = i
			indexpage = x
			q = x.item[i].ch
		case *dataPage:
			return false
		}
	}
}

func (t *Tree) overflow(indexPage *indexPage, dataPage *dataPage, pi, i int, key interface{}, value interface{}) {
	t.ver++
	leftDataPage, rightDataPage := indexPage.siblings(pi)
	if leftDataPage != nil && leftDataPage.count < 2*kd && i != 0 {
		leftDataPage.mvL(dataPage, 1)
		t.insert(dataPage, i-1, key, value)
		indexPage.item[pi-1].key = dataPage.item[0].key
		return
	}

	if rightDataPage != nil && rightDataPage.count < 2*kd {
		if i < 2*kd {
			dataPage.mvR(rightDataPage, 1)
			t.insert(dataPage, i, key, value)
			indexPage.item[pi].key = rightDataPage.item[0].key
			return
		}

		t.insert(rightDataPage, 0, key, value)
		indexPage.item[pi].key = key
		return
	}

	t.split(indexPage, dataPage, pi, i, key, value)
}

//TODO? and konw how to use
func (t *Tree) split(indexpage *indexPage, datapage *dataPage, pi, i int, key interface{}, value interface{}) {
	t.ver++
	r := btDataPagePoolEnity.Get().(*dataPage)
	//双向链表操作
	if datapage.next != nil {
		r.next = datapage.next
		r.next.prev = r
	} else {
		t.last = r
	}
	datapage.next = r
	r.prev = datapage

	copy(r.item[:], datapage.item[kd:2*kd])
	for i := range datapage.item[kd:] {
		datapage.item[kd+i] = zdataitem //gc
	}
	datapage.count = kd
	r.count = kd

	var done bool
	if i > kd {
		done = true
		t.insert(r, i-kd, key, value)
	}
	//更新index page
	if pi >= 0 {
		indexpage.insert(pi, r.item[0].key, r)
	} else {
		t.TreeR = newIndexPage(datapage).insert(0, r.item[0].key, r)
	}

	if done {
		return
	}

	t.insert(datapage, i, key, value)
}

//TODO ? and know how to use
func (t *Tree) splitIndexPage(indexpage *indexPage, indexpage2 *indexPage, pi, i int) (*indexPage, int) {
	t.ver++
	r := btIndexPagePoolEnity.Get().(*indexPage)
	copy(r.item[:], indexpage2.item[kx+1:])
	indexpage2.count = kx
	r.count = kx

	if pi >= 0 {
		indexpage.insert(pi, indexpage2.item[kx].key, r)
	} else {
		t.TreeR = newIndexPage(indexpage2).insert(0, indexpage2.item[kx].key, r)
	}

	indexpage2.item[kx].key = zkey
	for i := range indexpage2.item[kx+1:] {
		indexpage2.item[kx+i+1] = zindexitem
	}

	if i > kx {
		indexpage2 = r
		i -= kx + 1
	}
	return indexpage2, i
}

/**
由 index节点到 data节点映射
*/
func (t *Tree) underflow(indexpage *indexPage, datapage *dataPage, pi int) {
	t.ver++
	leftDataPage, rightDataPage := indexpage.siblings(pi)

	if leftDataPage != nil && leftDataPage.count+datapage.count >= 2*kd {
		leftDataPage.mvL(datapage, 1)
		datapage.item[pi].key = rightDataPage.item[0].key
		return
	}

	if rightDataPage != nil && datapage.count+rightDataPage.count >= 2*kd {
		datapage.mvL(rightDataPage, 1)
		indexpage.item[pi].key = rightDataPage.item[0].key
		rightDataPage.item[rightDataPage.count] = zdataitem //GC
		return
	}

	if leftDataPage != nil {
		t.cat(indexpage, leftDataPage, datapage, pi-1)
		return
	}

	t.cat(indexpage, datapage, rightDataPage, pi)
}

func (t *Tree) underflowIndexPage(indexpage *indexPage, indexpage2 *indexPage, pi, i int) (*indexPage, int) {
	t.ver++
	var left, right *indexPage
	if pi >= 0 {
		if pi > 0 {
			left = indexpage.item[pi-1].ch.(*indexPage)
		}
		if pi < indexpage.count {
			right = indexpage.item[pi+1].ch.(*indexPage)
		}
	}

	if left != nil && left.count > kx {
		indexpage2.item[indexpage2.count+1].ch = indexpage2.item[indexpage2.count].ch
		copy(indexpage2.item[1:], indexpage2.item[:indexpage2.count])
		indexpage2.item[0].ch = left.item[left.count].ch
		indexpage2.item[0].key = indexpage.item[pi-1].key
		indexpage2.count++
		i++
		left.count--
		indexpage.item[pi-1].key = left.item[left.count].key
		return indexpage2, i
	}

	if right != nil && right.count > kx {
		indexpage2.item[indexpage2.count].key = indexpage.item[pi].key
		indexpage2.count++
		indexpage2.item[indexpage2.count].ch = right.item[0].ch
		indexpage.item[pi].key = right.item[0].key
		copy(right.item[:], right.item[1:right.count])
		right.count--

		rightCount := right.count
		right.item[rightCount].ch = right.item[rightCount+1].ch
		right.item[rightCount].key = zkey
		right.item[rightCount].ch = nil
		return indexpage2, i
	}

	if left != nil {
		i += left.count + 1
		t.catIndexPage(indexpage, left, indexpage2, pi-1)
		indexpage2 = left
		return indexpage2, i
	}

	t.catIndexPage(indexpage, indexpage2, right, pi)
	return indexpage2, i
}

func (t *Tree) cat(p *indexPage, q, r *dataPage, pi int) {
	t.ver++
	q.mvL(r, r.count)
	if r.next != nil {
		r.next.prev = q
	} else {
		t.last = q
	}

	q.next = r.next
	*r = zdatapage
	btDataPagePoolEnity.Put(r)
	if p.count > 1 {
		p.extract(pi)
		p.item[pi].ch = q
		return
	}
	switch x := t.TreeR.(type) {
	case *indexPage:
		*x = zindexpage
		btIndexPagePoolEnity.Put(x)
	case *dataPage:
		*x = zdatapage
		btDataPagePoolEnity.Put(x)
	}

	t.TreeR = q
}

func (t *Tree) catIndexPage(indexpage, indexpage2, right *indexPage, pi int) {
	t.ver++
	indexpage2.item[indexpage2.count].key = indexpage.item[pi].key
	copy(indexpage2.item[indexpage2.count+1:], right.item[:right.count])
	indexpage2.count += right.count + 1
	indexpage2.item[indexpage2.count].ch = right.item[right.count].ch
	*right = zindexpage
	btIndexPagePoolEnity.Put(right)

	if indexpage.count > 1 {
		indexpage.count--
		indexpageCount := indexpage.count
		//中间插入
		if pi < indexpageCount {
			indexpage.item[pi].key = indexpage.item[pi+1].key
			copy(indexpage.item[pi+1:], indexpage.item[pi+2:indexpageCount+1])
			indexpage.item[indexpageCount].ch = indexpage.item[indexpageCount+1].ch
			indexpage.item[indexpageCount].key = zkey
			indexpage.item[indexpageCount+1].ch = nil
		}

		return
	}

	switch x := t.TreeR.(type) {
	case *indexPage:
		*x = zindexpage
		btIndexPagePoolEnity.Put(x)
	case *dataPage:
		*x = zdatapage
		btDataPagePoolEnity.Put(x)
	}

	t.TreeR = indexpage2
}

//----------------------------------Enumerator
//Close recycles e to a pool for possible later reuse. No referecnces to e
// should exist or such references must not be used afterwards.

func (e *Enumerator) Close() {
	*e = zenumerator
	btEnumeratorPoolEnity.Put(e)
}

// return the currently enumerated item, if it exists and moves to the
// next item in the key collation order. If there is no item to return, err == io.EOF is returned
func (e *Enumerator) Next() (key interface{}, value interface{}, err error) {
	if err = e.err; err != nil {
		return
	}

	if e.ver != e.tree.ver {
		f, _ := e.tree.Seek(e.key)
		*e = *f
		f.Close()
	}

	if e.enumeratorQ == nil {
		e.err, err = io.EOF, io.EOF
		return
	}

	if e.i >= e.enumeratorQ.count {
		if err = e.next(); err != nil {
			return
		}
	}

	item := e.enumeratorQ.item[e.i]
	key, value = item.key, item.value
	e.key, e.hit = key, true
	e.next()
	return
}

func (e *Enumerator) next() error {
	if e.enumeratorQ == nil {
		e.err = io.EOF
		return io.EOF
	}

	switch {
	case e.i < e.enumeratorQ.count-1:
		e.i++
	default:
		if e.enumeratorQ, e.i = e.enumeratorQ.next, 0; e.enumeratorQ == nil {
			e.err = io.EOF
		}
	}

	return e.err
}

// Prev returns the currently enumerated item, if it exists and moves to the
// previous item in the key collation order. If there is no item to return, err
// == io.EOF is returned.

func (e *Enumerator) Prev() (key interface{}, value interface{}, err error) {
	if err = e.err; err != nil {
		return
	}

	if e.ver != e.tree.ver {
		f, _ := e.tree.Seek(e.key)
		*e = *f
		f.Close()
	}

	if e.enumeratorQ == nil {
		e.err, err = io.EOF, io.EOF
		return
	}

	if !e.hit {
		// move to previous becasuse seek overshoots if there's no hit
		if err = e.prev(); err != nil {
			return
		}
	}

	if e.i > e.enumeratorQ.count {
		if err = e.prev(); err != nil {
			return
		}
	}

	item := e.enumeratorQ.item[e.i]
	key, value = item.key, item.value
	e.key, e.hit = key, true
	e.prev()
	return
}

func (e *Enumerator) prev() error {
	if e.enumeratorQ == nil {
		e.err = io.EOF
		return io.EOF
	}

	switch {
	case e.i > 0:
		e.i--
	default:
		if e.enumeratorQ = e.enumeratorQ.prev; e.enumeratorQ == nil {
			e.err = io.EOF
			break
		}

		e.i = e.enumeratorQ.count - 1

	}

	return e.err
}
