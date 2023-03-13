package dict

import (
	"sync"
)

type SyncDict struct {
	m sync.Map
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}

func (dict *SyncDict) Get(key string) (value any, exist bool) {
	value, exist = dict.m.Load(key)
	return
}
func (dict *SyncDict) Len() (result int) {
	dict.m.Range(func(key, value any) bool {
		result += 1
		return true
	})
	return
}
func (dict *SyncDict) Put(key string, value any) (result int) {
	dict.m.Store(key, value)
	return 1
}
func (dict *SyncDict) PutIfAbsent(key string, value any) (result int) {
	if _, ok := dict.m.Load(key); !ok {
		dict.m.Store(key, value)
		return 1
	}
	return

}
func (dict *SyncDict) PutIfExists(key string, value any) (result int) {
	if _, ok := dict.m.Load(key); ok {
		dict.m.Store(key, value)
		return 1
	}
	return
}

func (dict *SyncDict) Remove(key string) (result int) {
	if _, ok := dict.m.Load(key); ok {
		dict.m.Delete(key)
		return 1
	}
	return
}

func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.m.Range(func(key, value any) bool {
		result := consumer(key.(string), value)
		return result
	})
}

func (dict *SyncDict) Keys() (result []string) {
	dict.m.Range(func(key, value any) bool {
		result = append(result, key.(string))
		return true
	})
	return
}

func (dict *SyncDict) RandomKeys(limit int) (result []string) {
	for i := 1; i <= limit; i++ {
		dict.m.Range(func(key, value any) bool {
			result = append(result, key.(string))
			return false
		})
	}
	return
}

// 返回多个不重复的键
func (dict *SyncDict) RandomDistinctKeys(limit int) (result []string) {
	i := 1
	dict.m.Range(func(key, value any) bool {
		result = append(result, key.(string))
		i++
		return i <= limit
	})
	return
}

func (dict *SyncDict) Clear() {
	var newDict sync.Map
	dict.m = newDict
}
