package dict

import "sync"

type SyncDict struct {
	m sync.Map
}

func NewSyncDict() *SyncDict {
	return &SyncDict{}
}

func (dict *SyncDict) Get(key string) (val any, exists bool) {
	return dict.m.Load(key)
}

func (dict *SyncDict) Len() int {
	length := 0
	dict.m.Range(func(key, value any) bool {
		length++
		return true
	})
	return length
}

func (dict *SyncDict) Set(key string, val any) (rowAffected int) {
	_, existed := dict.m.Load(key)
	dict.m.Store(key, val)
	if existed {
		return 0
	}
	return 1
}

func (dict *SyncDict) SetIfAbsent(key string, val any) (rowAffected int) {
	if _, existed := dict.m.Load(key); existed {
		return 0
	}
	dict.m.Store(key, val)
	return 1
}

func (dict *SyncDict) SetIfExists(key string, val any) (rowAffected int) {
	if _, existed := dict.m.Load(key); existed {
		dict.m.Store(key, val)
		return 1
	}
	return 0
}

func (dict *SyncDict) Remove(key string) (rowAffected int) {
	if _, existed := dict.m.Load(key); existed {
		dict.m.Delete(key)
		return 1
	}
	return 0
}

func (dict *SyncDict) Range(consumer func(key string, val any) bool) {
	dict.m.Range(func(key, value any) bool {
		consumer(key.(string), value)
		return true
	})
}

func (dict *SyncDict) Keys() []string {
	result := make([]string, dict.Len())
	dict.m.Range(func(key, value any) bool {
		result = append(result, key.(string))
		return true
	})
	return result
}

func (dict *SyncDict) RandomKeys(limit int) []string {
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		dict.m.Range(func(key, value any) bool {
			result[i] = key.(string)
			return false
		})
	}
	return result
}

func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	result := make([]string, limit)
	i := 0
	dict.m.Range(func(key, value any) bool {
		result[i] = key.(string)
		i++
		if i == limit {
			return false
		}
		return true
	})
	return result
}

func (dict *SyncDict) Clear() {
	*dict = *NewSyncDict()
}
