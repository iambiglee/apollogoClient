package memory

import (
	"github.com/apollogoClient/v1/agache"
	"sync"
)

// DefaultCacheFactory TODO 空结构体为什么这里面什么都没有
type DefaultCacheFactory struct {
}

// DefaultCache  为什么结构体可以只写类型不写引用，没有默认用数据类型代替名字
type DefaultCache struct {
	defaultCache sync.Map
	count        int64
}

func (d DefaultCache) Set(key string, value interface{}, expireSeconds int) (err error) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultCache) EntryCount() (entryCount int64) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultCache) Get(key string) (value interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultCache) Del(key string) (affect bool) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultCache) Range(f func(key interface{}, value interface{}) bool) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultCache) Clear() {
	//TODO implement me
	panic("implement me")
}

// Create TODO 为什么会自动实现Create方法
func (d DefaultCacheFactory) Create() agache.CacheInterface {
	return &DefaultCache{}
}
