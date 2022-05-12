package storage

import (
	"container/list"
	"github.com/apollogoClient/v1/agache"
	"github.com/apollogoClient/v1/env/config"
	"github.com/apollogoClient/v1/extension"
	"sync"
	"sync/atomic"
)

type Config struct {
	namespace string
	cache     agache.CacheInterface
	isInit    atomic.Value
	waitInit  sync.WaitGroup
}

type Cache struct {
	apolloConfigCache sync.Map
	ChangeListener    *list.List
}

// CreateNamespaceConfig 根据namespace初始化goClient 内部配置
//SplitNamespaces() 是个什么原理：利用参数就是一个接口的方式，直接通过方法实现方法
func CreateNamespaceConfig(namespace string) *Cache {
	var apolloConfigCache sync.Map
	config.SplitNamespaces(namespace, func(namespace string) {
		if _, ok := apolloConfigCache.Load(namespace); ok {
			return
		}
		c := initConfig(namespace, extension.GetCacheFactory())
		apolloConfigCache.Store(namespace, c)
	})
	return &Cache{
		apolloConfigCache: apolloConfigCache,
		ChangeListener:    list.New(),
	}
}

func initConfig(namespace string, factory agache.CacheFactory) *Config {
	c := &Config{
		namespace: namespace,
		cache:     factory.Create(),
	}
	c.isInit.Store(false)
	c.waitInit.Add(1)
	return c
}
