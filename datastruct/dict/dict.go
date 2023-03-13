package dict

type Consumer func(key string, value any) bool

type Dict interface {
	Get(key string) (value any, exist bool)
	Len() (result int)
	Put(key string, value any) (result int)
	PutIfAbsent(key string, value any) (result int)
	PutIfExists(key string, value any) (result int)
	Remove(key string) (result int)
	ForEach(consumer Consumer)
	Keys() (result []string)
	RandomKeys(limit int) (result []string)
	//返回多个不重复的键
	RandomDistinctKeys(limit int) (result []string)
	Clear()
}
