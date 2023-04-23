package list

type Link struct {
	first *node
	last  *node
	len   int
}

type node struct {
	next  *node
	prev  *node
	value interface{}
}

func (l *Link) Add(val interface{}) {
	if l == nil {
		panic("Link is nil")
	}
	n := &node{
		value: val,
	}
	if l.last == nil {
		//empty list
		l.first = n
		l.last = n
	} else {
		n.prev = l.last
		l.last.next = n
		l.last = n
	}
	l.len++
}

func (l *Link) find(index int) (n *node) {
	if index < l.len/2 {
		for i := 0; i < index; i++ {
			n = n.next
		}
	} else {
		n := l.last
		for i := l.len - 1; i > index; i-- {
			n = n.prev
		}
	}
	return n
}

func (l *Link) Get(index int) (val interface{}) {
	if l == nil {
		panic("Link is nil")
	}
	if index < 0 || index >= l.len {
		panic("index out of range")
	}
	return l.find(index).value
}

func (l *Link) Set(index int, val interface{}) {
	if l == nil {
		panic("Link is nil")
	}
	if index < 0 || index >= l.len {
		panic("index out of range")
	}
	n := l.find(index)
	n.value = val
}

func (l *Link) Insert(index int, val interface{}) {
	if l == nil {
		panic("Link is nil")
	}
	if index < 0 || index >= l.len {
		panic("index out of range")
	}
	if index == l.len {
		l.Add(val)
		return
	}
	n := l.find(index)
	newNode := &node{
		value: val,
	}
	if n.prev == nil {
		l.first = newNode
	} else {
		prevNode := n.prev
		prevNode.next = newNode
		newNode.prev = prevNode
	}
	n.prev = newNode
	newNode.next = n
	l.len++
}

func (l *Link) Remove(index int) (val interface{}) {
	if l == nil {
		panic("Link is nil")
	}
	if index < 0 || index >= l.len {
		panic("index out of range")
	}
	rmNode := l.find(index)
	val = rmNode.value
	l.delete(rmNode)
	return
}

func (l *Link) RemoveLast() (val interface{}) {
	if l == nil {
		panic("Link is nil")
	}
	if l.last == nil {
		return nil
	}
	val = l.last.value
	l.delete(l.last)
	return
}

func (l *Link) Len() int {
	if l == nil {
		panic("Link is nil")
	}
	return l.len
}

func (l *Link) delete(node *node) {
	if node.next == nil {
		//删除最后一个
		l.last = node.prev
	} else {
		node.next.prev = node.prev
	}
	if node.prev == nil {
		//删除第一个
		l.first = node.next
	} else {
		node.prev.next = node.next
	}
	node.next = nil
	node.prev = nil
	l.len--
}

// RemoveByVal 移除最多指定count数量的值
// 从左向右扫描
func (l *Link) RemoveByVal(expected Expected, count int) (reCount int) {
	if l == nil {
		panic("Link is nil")
	}
	if l.first == nil {
		return
	}
	for n := l.first; n.next != nil; n = n.next {
		if expected(n.value) {
			l.delete(n)
			reCount++
			if reCount == count {
				break
			}
		}
	}
	return
}

// ReverseRemoveByVal 移除最多指定count数量的值
// 从右向左边扫描
func (l *Link) ReverseRemoveByVal(expected Expected, count int) (reCount int) {
	if l == nil {
		panic("Link is nil")
	}
	if l.last == nil {
		return
	}
	for n := l.last; n.prev != nil; n = n.prev {
		if expected(n.value) {
			l.delete(n)
			reCount++
			if reCount == count {
				break
			}
		}
	}
	return
}

func (l *Link) RemoveAllByVal(expected Expected) (rmCount int) {
	if l == nil {
		panic("Link is nil")
	}
	if l.first == nil {
		return
	}
	for n := l.first; n.next != nil; n = n.next {
		if expected(n.value) {
			l.delete(n)
			rmCount++
		}
	}
	return
}

func (l *Link) ForEach(consumer Consumer) {
	if l == nil {
		panic("Link is nil")
	}
	n := l.first
	i := 0
	for n != nil {
		goNext := consumer(i, n.value)
		if !goNext {
			break
		}
		n = n.next
		i++
	}
}

// 返回指定的值是否存在列表中
func (l *Link) Contains(expected Expected) (contains bool) {
	if l == nil {
		panic("Link is nil")
	}
	l.ForEach(func(index int, val interface{}) bool {
		if expected(val) {
			contains = true
			return true
		}
		return false
	})
	return
}

func (l *Link) Range(start int, end int) []interface{} {
	if l == nil {
		panic("Link is nil")
	}
	if start > l.len || start < 0 {
		panic("start out of range")
	}
	if end > l.len || end < start || end < 0 {
		panic("end out of range")
	}
	vals := make([]interface{}, end-start)
	node := l.find(start)
	index := 0
	i := 0
	for node != nil {
		if i >= start && i < end {
			vals[index] = node.value
			index++
		} else if i >= end {
			break
		}
		i++
		node = node.next
	}
	return vals
}

func Make(vals ...interface{}) *Link {
	Link := &Link{}
	for _, val := range vals {
		Link.Add(val)
	}
	return Link
}
