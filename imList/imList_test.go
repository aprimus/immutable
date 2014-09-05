package imList

import "testing"

func createTestList() *IMList {
	l := New()
	for i := 9; i >= 0; i-- {
		l = l.Push(i)
	}
	return l
}

func makeSpotInt(target int) func(interface{}) bool {
	return func(v interface{}) bool {
		vi := v.(int)
		if vi == target {
			return true
		}
		return false
	}
}

func confirmListMatch(l *IMList, targets []int, t *testing.T) *testing.T {
	i := 0
	for e := range l.Iter() {
		v := e.Value.(int)
		if v != targets[i] {
			t.Error("Unexpected return from the list. Expected/Return:", targets[i], "/", v)
		}
		i++
	}
	return t
}

func TestRemove(t *testing.T) {
	l := createTestList()
	l2 := l.Remove(2)
	desiredState := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	i := 0
	// Check that original list was not altered
	for e := range l.Iter() {
		v := e.Value.(int)
		//println("Expecting", desiredState[i], "got", v)
		if v != desiredState[i] {
			t.Error("Unexpected return from the list, possibly not immutable. Expected/Return:", desiredState[i], "/", v)
		}
		i++
	}
	// Test that derivative list was altered
	desiredState = []int{0, 1, 3, 4, 5, 6, 7, 8, 9}
	i = 0
	for e := range l2.Iter() {
		v := e.Value.(int)
		//println("Expecting", desiredState[i], "got", v)
		if v != desiredState[i] {
			t.Error("Unexpected return from the list. Expected/Return:", desiredState[i], "/", v)
		}
		i++
	}

}

func TestUpdateOrInsert(t *testing.T) {
	l := createTestList()
	l2 := l.UpdateOrInsert(12, makeSpotInt(2))
	l3 := l2.UpdateOrInsert(10, makeSpotInt(10))
	desiredState := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	i := 0
	// Test that original list unaltered
	for e := range l.Iter() {
		v := e.Value.(int)
		if v != desiredState[i] {
			t.Error("Unexpected return from the list. Expected/Return:", desiredState[i], "/", v)
		}
		i++
	}
	if i != len(desiredState) {
		t.Error("Original list only had", i-1, "elements. Failed to match", desiredState[i])
	}
	// Test that l2 has changed 2 -> 12
	desiredState = []int{0, 1, 12, 3, 4, 5, 6, 7, 8, 9}
	i = 0
	for e := range l2.Iter() {
		v := e.Value.(int)
		if v != desiredState[i] {
			t.Error("Unexpected return from the list. Expected/Return:", desiredState[i], "/", v)
		}
		i++
	}
	if i != len(desiredState) {
		t.Error("Altered List only had", i, "elements")
	}
	// Change that l3 has the value "10" added
	desiredState = []int{10, 0, 1, 12, 3, 4, 5, 6, 7, 8, 9}
	i = 0
	for e := range l3.Iter() {
		v := e.Value.(int)
		if v != desiredState[i] {
			t.Error("Unexpected return from the list. Expected/Return:", desiredState[i], "/", v)
		}
		i++
	}
	if i != len(desiredState) {
		t.Error("List (after insertion) only had", i, "elements")
	}
}

func trueIfEven(i interface{}) bool {
	if i.(int)%2 == 0 {
		return true
	}
	return false
}

func TestIterFilter(t *testing.T) {
	l := createTestList()
	ch := l.IterFilter(trueIfEven)
	desiredState := []int{0, 2, 4, 6, 8}
	i := 0
	for e := range ch {
		v := e.(int)
		if v != desiredState[i] {
			t.Error("Unexpected return from the list. Expected/Return:", desiredState[i], "/", v)
		}
		i++
	}
	if i != len(desiredState) {
		t.Error("List only had", i, "elements")
	}
}

func TestRemoveByFunc(t *testing.T) {
	l := createTestList()
	l1 := l.RemoveByFunc(makeSpotInt(1))
	desiredState := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	i := 0
	for e := range l.Iter() {
		v := e.Value.(int)
		if v != desiredState[i] {
			t.Error("Unexpected return from the list. Expected/Return:", desiredState[i], "/", v)
		}
		i++
	}
	if i != len(desiredState) {
		t.Error("List only had", i-1, "elements. Failed to match", desiredState[i])
	}

	desiredState = []int{0, 2, 3, 4, 5, 6, 7, 8, 9}
	i = 0
	for e := range l1.Iter() {
		v := e.Value.(int)
		if v != desiredState[i] {
			t.Error("Unexpected return from the list. Expected/Return:", desiredState[i], "/", v)
		}
		i++
	}
	if i != len(desiredState) {
		t.Error("List only had", i-1, "elements. Failed to match", desiredState[i])
	}
}
