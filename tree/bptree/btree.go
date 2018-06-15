package bptree

// import (
// 	"fmt"
// 	"sync"
// )

// const (
// 	kx = 32 //TODO benchmark tune this number if using custom key/value type(s)
// 	kd = 32 //TODO benchmark tune this number if using custom key/value type(s)
// )

// func init() {
// 	if kd < 1 {
// 		panic(fmt.Errorf("kd %d: out of range", kd))
// 	}

// 	if kx < 2 {
// 		panic(fmt.Errorf("kx %d: out of range", kx))
// 	}
// }

// type btEpool struct{ sync.Pool }

// func (p *btEpool) get(err error, hit bool, i int, k interface{}, q *d, t *Tree, ver int64) *Enumerator {
// 	x := p.Get().(*Enumerator)
// 	x.err, x.hit, x.i, x.k, x.q, x.t, x.ver = err, hit, i, k, q, t, ver
// 	return x
// }

// type btTpool struct {
// 	sync.Pool
// }

// func (p *btTpool) get(cmp Cmp) *Tree {
// 	x := p.Get().(*Tree)
// 	x.cmp = cmp
// 	return x
// }

// var (
// 	btDPool = sync.Pool{New: func() interface{} { return &d{} }}
// 	btEPool = btEpool{sync.Pool{New: func() interface{} { return &Enumerator{} }}}
// 	btTPool = btTpool{sync.Pool{New: func() interface{} { return &Tree{} }}}
// 	btXPool = sync.Pool{New: func() interface{} { return &x{} }}
// )

// type (
// 	// Cmp compares a and b. Return value is:
// 	//
// 	//	< 0 if a <  b
// 	//	  0 if a == b
// 	//	> 0 if a >  b
// 	//
// 	Cmp func(a, b interface{}) int

// 	d struct { //数据页
// 		c int
// 		d [2*kd + 1]de
// 		n *d //next
// 		p *d //prev
// 	}

// 	de struct { //数据 item
// 		k interface{} //key
// 		v interface{} //value
// 	}

// 	/**
// 	Enumerator captures the state of enumerating a tree. It is returned
// 	from the Seek* methods. The enumerator is aware of any mutations
// 	made to the tree in the process of enumerationg it andautomatically
// 	resumes the enumeration at the proper key, if possible.

// 	However, once an Enumerator returns io.EOF to signal "no more
// 	items", it does no more attempt to "resync" on tree mutation. In
// 	other words, io.EOF form an Enumerator is "sticky"(idemotent).
// 	*/
// 	Enumerator struct {
// 		err error
// 		hit bool
// 		i   int
// 		k   interface{}
// 		q   *d
// 		t   *Tree
// 		ver int64
// 	}

// 	Tree struct { //B+treee
// 		c     int
// 		cmp   Cmp
// 		first *d
// 		last  *d
// 		r     interface{}
// 		ver   int64
// 	}

// 	xe struct { //索引条目
// 		ch interface{}
// 		k  interface{}
// 	}

// 	x struct { //索引页面
// 		c int
// 		x [2*kx + 2]xe
// 	}
// )

// var (
// 	zd  d
// 	zde de
// 	ze  Enumerator
// 	zk  interface{}
// 	zt  Tree
// 	zx  x
// 	zxe xe
// )

// func clr(q interface{}) {
// 	switch x := q.(type) {
// 	case *x:
// 		for i := 0; i < x.c; i++ {
// 			clr(x.x[i].ch)
// 		}
// 		*x = zx
// 		btXPool.Put(x)
// 	case *d:
// 		*x = zd
// 		btDPool.Put(x)
// 	}
// }

// //x
// func newX(ch0 interface{}) *x {
// 	r := btXPool.Get().(*x)
// 	r.x[0].ch = ch0

// 	return r
// }

// /**
// 提取 x中第i 个xe
// */
// func (q *x) extract(i int) {
// 	q.c--
// 	if i < q.c {
// 		copy(q.x[i:], q.x[i+1:q.c+1])
// 		q.x[q.c].ch = q.x[q.c+1].ch
// 		q.x[q.c].k = zk  //GC
// 		q.x[q.c+1] = zxe //GC
// 	}
// }

// func (q *x) insert(i int, k interface{}, ch interface{}) *x {
// 	c := q.c
// 	if i < c {
// 		q.x[c+1].ch = q.x[c].ch
// 		copy(q.x[i+2:], q.x[i+1:c])
// 		q.x[i+1].k = q.x[i].k
// 	}

// 	c++
// 	q.c = c
// 	q.x[i].k = k
// 	q.x[i+1].ch = ch
// 	return q
// }

// /**
// 兄弟元素
// */
// func (q *x) siblings(i int) (l, r *d) {
// 	if i >= 0 {
// 		l = q.x[i-1].ch.(*d)
// 	}
// 	if i < q.c {
// 		r = q.x[i+1].ch.(*d)
// 	}

// 	return
// }

// //--------------------------------d
// //TODO understand there mvL and mvR
// func (l *d) mvL(r *d, c int) {
// 	copy(l.d[l.c:], r.d[:c])
// 	copy(r.d[:], r.d[c:r.c])
// 	l.c += c
// 	r.c -= c
// }

// func (l *d) mvR(r *d, c int) {
// 	copy(r.d[c:], r.d[:r.c])
// 	copy(r.d[:c], l.d[l.c-c:])
// 	l.c -= c
// 	r.c += c
// }

// //--------------------------------tree
// //TreeNew returns a newly created, empty Tree. The compare function is used
// //for key collation
// func TreeNew(cmp Cmp) *Tree {
// 	return btTPool.get(cmp)
// }

// func (t *Tree) Clear() {
// 	if t.r == nil {
// 		return
// 	}

// 	clr(t.r)
// 	t.c, t.first, t.last, t.r = 0, nil, nil, nil
// 	t.ver++
// }

// /**
// close performs clear and recycles t to a pool for possible later reuse. No
// references to t should exits or such references must not  be used afterwards.
// */
// func (t *Tree) Close() {
// 	t.Clear()
// 	*t = zt
// 	btTPool.Put(t)
// }

// func (t *Tree) cat(p *x, q, r *d, pi int) {
// 	t.ver++
// 	q.mvL(r, r.c)
// 	if r.n != nil {
// 		r.n.p = q
// 	} else {
// 		t.last = q
// 	}
// 	q.n = r.n
// 	*r = zd
// 	btDPool.Put(r)
// 	if p.c > 1 {
// 		p.extract(pi)
// 		p.x[pi].ch = q
// 		return
// 	}

// 	switch x := t.r.(type) {
// 	case *x:
// 		*x = zx
// 		btXPool.Put(x)
// 	case *d:
// 		*x = zd
// 		btDPool.Put(x)
// 	}

// 	t.r = q
// }

// func (t *Tree) catX(p, q, r *x, pi int) {
// 	t.ver++
// 	q.x[q.c].k = p.x[pi].k
// 	copy(q.x[q.c+1:], r.x[:r.c])
// 	q.c += r.c + 1
// 	q.x[q.c].ch = r.x[r.c].ch
// 	*r = zx
// 	btXPool.Put(r)
// 	if p.c > 1 {
// 		p.c--
// 		pc := p.c
// 		if pi < pc {
// 			p.x[pi].k = p.x[pi+1].k
// 			//todo
// 		}
// 	}
// }

// func (t *Tree) insert(q *d, i int, k interface{}, v interface{}) *d {
// 	t.ver++
// 	c := q.c
// 	if i < c {
// 		copy(q.d[i+1:], q.d[i:c])
// 	}
// 	c++
// 	q.c = c
// 	q.d[i].k, q.d[i].v = k, v
// 	t.c++
// 	return q
// }

// func (t *Tree) find(q interface{}, k interface{}) (i int, ok bool) {
// 	var mk interface{}
// 	l := 0
// 	switch x := q.(type) {
// 	case *x:
// 		h := x.c - 1
// 		for l <= h {
// 			m := (l + h) >> 1
// 			mk = x.x[m].k
// 			switch cmp := t.cmp(k, mk); {
// 			case cmp > 0:
// 				l = m + 1
// 			case cmp == 0:
// 				return m, true
// 			default:
// 				h = m - 1
// 			}
// 		}
// 	case *d:
// 		h := x.c - 1
// 		for l <= h {
// 			m := (l + h) >> 1
// 			mk = x.d[m].k
// 			switch cmp := t.cmp(k, mk); {
// 			case cmp > 0:
// 				l = m + 1
// 			case cmp == 0:
// 				return m, true
// 			default:
// 				h = m - 1
// 			}
// 		}
// 	}
// 	return l, false
// }

// func (t *Tree) extract(q *d, i int) {
// 	t.ver++
// 	q.c--
// 	if i < q.c {
// 		copy(q.d[i:], q.d[i+1:q.c+1])
// 	}
// 	q.d[q.c] = zde //gc
// 	t.c--
// }

// func (t *Tree) Delete(k interface{}) (ok bool) {
// 	pi := -1
// 	var p *x
// 	q := t.r
// 	if q == nil {
// 		return false
// 	}
// 	for {
// 		var i int
// 		i, ok = t.find(q, k)
// 		if ok {
// 			switch x := q.(type) {
// 			case *x:
// 				if x.c < kx && q != t.r {
// 					x, i = t.unserflowX(p, x, pi, i)
// 				}
// 				pi = i + 1
// 				p = x
// 				q = x.x[pi].ch
// 				continue
// 			case *d:
// 				t.extract(x, i)
// 				if x.c >= kd {
// 					return true
// 				}
// 				if q != t.r {
// 					t.unserflow(p, x, pi)
// 				} else if t.c == 0 {
// 					t.Clear()
// 				}
// 				return true
// 			}
// 		}
// 		switch x := q.(type) {
// 		case *x:
// 			if x.c < kx && q != t.r {
// 				x, i = t.unserflowX(p, x, pi, i)
// 			}
// 			pi = i
// 			p = x
// 			q = x.x[i].ch
// 		case *d:
// 			return false
// 		}
// 	}

// }

// //first returns the first item of the tree in the key collating order,
// //or (zero-value, zero-value) if the tree is empty
// func (t *Tree) First() (k interface{}, v interface{}) {
// 	if q := t.first; q != nil {
// 		q := &q.d[0]
// 		k, v = q.k, q.v
// 	}

// 	return
// }

// //return the last tiem of the tree in thre key collating order, or
// //(zero-value, zero-value) if the tree is empty
// func (t *Tree) Last() (k interface{}, v interface{}) {
// 	if q := t.last; q != nil {
// 		q := &q.d[q.c-1]
// 		k, v = q.k, q.v
// 	}
// 	return
// }

// //Get returns the value associated with k and true if it exists. Otherwise Get
// // returns (zero-value, false)
// func (t *Tree) Get(k interface{}) (v interface{}, ok bool) {
// 	q := t.r
// 	if q == nil {
// 		return
// 	}

// 	for {
// 		var i int
// 		if i, ok = t.find(q, k); ok {
// 			switch x := q.(type) {
// 			case *x:
// 				q = x.x[i+1].ch
// 				continue
// 			case *d:
// 				return x.d[i].v, true
// 			}
// 		}

// 		switch x := q.(type) {
// 		case *x:
// 			q = x.x[i].ch
// 		default:
// 			return
// 		}
// 	}
// }

// // return the number of items in the tree
// func (t *Tree) Len() int {
// 	return t.c
// }

// func (tree *Tree) overflow(indexPage *x, dataPage *d, pi, i int, key interface{}, value interface{}) {
// 	tree.ver++
// 	left, right := indexPage.siblings(pi)
// 	if left != nil && left.c < 2*kd && i != 0 {
// 		left.mvL(dataPage, 1)
// 		tree.insert(dataPage, i-1, key, value)
// 		indexPage.x
// 	}
// }
