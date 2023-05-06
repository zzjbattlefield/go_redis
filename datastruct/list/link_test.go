package list

import (
	"strconv"
	"testing"
)

func TestAdd(t *testing.T) {
	list := Make()
	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	list.ForEach(func(i int, val interface{}) bool {
		intval, _ := val.(int)
		if intval != i {
			t.Errorf("Add error, expected %d, but got %d", i, intval)
		}
		return true
	})
}

func TestGet(t *testing.T) {
	list := Make()
	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	for i := 0; i < 10; i++ {
		val := list.Get(i)
		intval, _ := val.(int)
		if intval != i {
			t.Errorf("Get error, expected %d, but got %d", i, intval)
		}
	}
}

func TestInsert(t *testing.T) {
	list := Make()
	for i := 0; i < 10; i++ {
		list.Insert(i, i)
	}
	list.ForEach(func(i int, val interface{}) bool {
		intval, _ := val.(int)
		if intval != i {
			t.Errorf("Insert error, expected %d, but got %d", i, intval)
		}
		return true
	})
}

func TestRemove(t *testing.T) {
	list := Make()
	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	for i := 9; i >= 0; i-- {
		list.Remove(i)
		if i != list.Len() {
			t.Errorf("Remove error, expected %d, but got %d", i, list.Len())
		}
		list.ForEach(func(i int, val interface{}) bool {
			intval, _ := val.(int)
			if i != intval {
				t.Errorf("Remove error, expected %d, but got %d", i, intval)
			}
			return true
		})
	}
}

func TestRemoveVal(t *testing.T) {
	list := Make()
	for i := 0; i < 10; i++ {
		list.Add(i)
		list.Add(i)
	}
	for index := 0; index < list.Len(); index++ {
		list.RemoveAllByVal(func(a interface{}) bool {
			return Equals(a, index)
		})
		list.ForEach(func(i int, v interface{}) bool {
			intVal, _ := v.(int)
			if intVal == index {
				t.Error("remove test fail: found  " + strconv.Itoa(index) + " at index: " + strconv.Itoa(i))
			}
			return true
		})
	}

	list = Make()
	for i := 0; i < 10; i++ {
		list.Add(i)
		list.Add(i)
	}
	for index := 0; index < list.Len(); index++ {
		list.RemoveByVal(func(a interface{}) bool {
			return Equals(a, index)
		}, 1)
	}
	list.ForEach(func(i int, v interface{}) bool {
		intVal, _ := v.(int)
		if intVal != i {
			t.Error("remove test fail: found  " + strconv.Itoa(i) + " at index: " + strconv.Itoa(i))
		}
		return true
	})

	list = Make()
	for i := 0; i < 10; i++ {
		list.Add(i)
		list.Add(i)
	}
	for i := 0; i < 10; i++ {
		list.ReverseRemoveByVal(func(a interface{}) bool {
			return a == i
		}, 1)
	}
	list.ForEach(func(i int, v interface{}) bool {
		intVal, _ := v.(int)
		if intVal != i {
			t.Error("test fail: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal))
		}
		return true
	})
	for i := 0; i < 10; i++ {
		list.ReverseRemoveByVal(func(a interface{}) bool {
			return a == i
		}, 1)
	}
	if list.Len() != 0 {
		t.Error("test fail: expected 0, actual: " + strconv.Itoa(list.Len()))
	}
}

func TestRange(t *testing.T) {
	list := Make()
	size := 10
	for i := 0; i < size; i++ {
		list.Add(i)
	}
	for start := 0; start < size; start++ {
		for end := start; end < size; end++ {
			slice := list.Range(start, end)
			if len(slice) != end-start {
				t.Error("expected " + strconv.Itoa(end-start) + ", get: " + strconv.Itoa(len(slice)) +
					", range: [" + strconv.Itoa(start) + "," + strconv.Itoa(end) + "]")
			}
			sliceIndex := 0
			for i := start; i < end; i++ {
				val := slice[sliceIndex]
				intval, _ := val.(int)
				if intval != i {
					t.Error("expected " + strconv.Itoa(i) + ", get: " + strconv.Itoa(intval) +
						", range: [" + strconv.Itoa(start) + "," + strconv.Itoa(end) + "]")
				}
				sliceIndex++
			}
		}

	}
}

func Equals(a interface{}, b interface{}) bool {
	sliceA, okA := a.([]byte)
	sliceB, okB := b.([]byte)
	if okA && okB {
		return BytesEquals(sliceA, sliceB)
	}
	return a == b
}

// BytesEquals check whether the given bytes is equal
func BytesEquals(a []byte, b []byte) bool {
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	size := len(a)
	for i := 0; i < size; i++ {
		av := a[i]
		bv := b[i]
		if av != bv {
			return false
		}
	}
	return true
}
