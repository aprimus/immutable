package imList

type IMList struct {
	head *node
}

type node struct {
	payload interface{}
	next    *node
}

type Element struct {
	Value interface{}
	Index int
}

func New() *IMList {
	return &IMList{nil}
}

func (iml *IMList) Push(payload interface{}) *IMList {
	n := &node{payload, iml.head}
	return &IMList{n}
}

func (iml *IMList) Pop() (interface{}, *IMList) {
	if iml.head == nil {
		return nil, &IMList{nil}
	}
	niml := &IMList{iml.head.next}
	return iml.head.payload, niml
}

func (iml *IMList) Iter() chan Element {
	out := make(chan Element)
	go func() {
		for n, i := iml.head, 0; n != nil; n, i = n.next, i+1 {
			out <- Element{n.payload, i}
		}
		close(out)
	}()
	return out
}

func (iml *IMList) Remove(index int) *IMList {
	niml := &IMList{nil}
	if iml.head == nil {
		return niml
	}
	nn := &node{nil, nil}
	if index == 0 {
		niml.head = iml.head.next
		return niml
	} else {
		nn = &node{iml.head.payload, nil}
		niml.head = nn
	}
	for n, i := iml.head.next, 1; n != nil; n, i = n.next, i+1 {
		if i != index {
			toAdd := &node{n.payload, nil}
			nn.next = toAdd
			nn = toAdd
		}
	}
	return niml
}

func (iml *IMList) RemoveByFunc(f func(interface{}) bool) *IMList {
	if iml.head == nil {
		return &IMList{nil}
	}
	if f(iml.head.payload) == true {
		return &IMList{iml.head.next}
	}
	newHead := &node{iml.head.payload, nil}
	niml := &IMList{newHead}
	prior := newHead
	for n := iml.head.next; n != nil; n = n.next {
		if f(n.payload) != true {
			thisNode := &node{n.payload, nil}
			prior.next = thisNode
			prior = thisNode
		} else { // We have a match.  Point to rest of list
			prior.next = n.next
			return niml
		}
	}
	// If we get to here, it means we never had a match,
	// so there's no point in making another copy of the
	// list, we can just return the old one.
	return iml
}

func (iml *IMList) Fetch(f func(interface{}) bool) interface{} {
	for n := iml.head; n != nil; n = n.next {
		if f(n.payload) {
			return n.payload
		}
	}
	return nil
}

func (iml *IMList) IterFilter(f func(interface{}) bool) chan interface{} {
	out := make(chan interface{})
	go func() {
		for n := iml.head; n != nil; n = n.next {
			if f(n.payload) {
				out <- n.payload
			}
		}
		close(out)
	}()
	return out
}

func (iml *IMList) UpdateOrInsert(payload interface{},
	f func(interface{}) bool) *IMList {
	niml := &IMList{nil}
	old := iml.Fetch(f)
	if old == nil { // Can just push at head
		newNode := &node{payload, iml.head}
		niml.head = newNode
		return niml
	}
	nodeToAdd := &node{payload, nil}
	niml.head = nodeToAdd
	n := iml.head
	for ; n != nil; n = n.next {
		if f(n.payload) == false {
			nextNode := &node{nil, nil}
			nodeToAdd.payload = n.payload
			nodeToAdd.next = nextNode
			nodeToAdd = nextNode
		} else {
			nodeToAdd.payload = payload
			nodeToAdd.next = n.next
			return niml
		}
	}
	panic("Should never get here -- either it was not found by fetch, or it was...")
}
