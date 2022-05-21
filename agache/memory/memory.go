package memory

import (
	"github.com/apollogoClient/v1/agache"
	"sync"
)

// DefaultCacheFactory  空结构体为什么这里面什么都没有
//空结构体是GO语言中的一种使用，情况使用空结构体的考量其实是：根本不关心结构体变量的值。
//适用于实现Set结构、在 channel 中传输信号,实现方法接收者,实现集合类型，need time to do it
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

// Create 为什么会自动实现Create方法，因为实现对象的时候，必须要一个实例对象
func (d DefaultCacheFactory) Create() agache.CacheInterface {
	return &DefaultCache{}
}
