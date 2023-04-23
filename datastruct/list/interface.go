package list

type Consumer func(i int, val interface{}) bool

type Expected func(val interface{}) bool

type List interface {
	Add(val interface{})
	Get(index int) (val interface{})
	Set(index int, val interface{})
	Insert(index int, val interface{})
	Remove(index int) (val interface{})
	RemoveLast() (val interface{})
	RemoveAllByVal(expected Expected) (rmCount int)
	RemoveByVal(expected Expected, count int) (reCount int)
	ReverseRemoveByVal(expected Expected, count int) (reCount int)
	Len() int
	ForEach(consumer Consumer)
	Contains(expected Expected) (contains bool)
	Range(start int, stop int) []interface{}
}
