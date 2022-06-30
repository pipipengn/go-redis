package dict

type Interface interface {
	Get(key string) (val any, exists bool)
	Len() int
	Set(key string, val any) (rowAffected int)
	SetIfAbsent(key string, val any) (rowAffected int)
	SetIfExists(key string, val any) (rowAffected int)
	Remove(key string) (rowAffected int)
	Range(consumer func(key string, val any) bool)
	Keys() []string
	RandomKeys(limit int) []string
	RandomDistinctKeys(limit int) []string
	Clear()
}

type DataEntity struct {
	Data any
}
