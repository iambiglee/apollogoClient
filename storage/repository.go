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

const (
	//1 minute
	configCacheExpireTime = 120

	defaultNamespace = "application"

	propertiesFormat = "%s=%v\n"
)

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

// UpdateApolloConfig config.Appconfig为什么这不能用*，为什么这里要用方法
//根据 config server 返回的内容更新并判断是否要写备份文件
func (c *Cache) UpdateApolloConfig(apolloConfig *config.ApolloConfig, appConfigFunc func() config.AppConfig) {
	if apolloConfig == nil {
		return
	}
	appConfig := appConfigFunc()
	appConfig.SetCurrentApolloConfig(&apolloConfig.ApolloConnConfig)
	c.UpdateApolloConfigCache(apolloConfig.Configurations, configCacheExpireTime, apolloConfig.NamespaceName)
}

//UpdateApolloConfigCache 根据conf[ig server返回的内容更新内存
func (c *Cache) UpdateApolloConfigCache(configurations map[string]interface{}, time int, namespace string) {
	c.GetConfig(namespace)
}

// GetConfig 根据namespace 获取Apollo配置
// 为什么* 在里面，为什么指定的Cache前面要有一个*,因为方法也是值传递，可以吧原来的cache 对象传递过来
func (c *Cache) GetConfig(namespace string) *Config {
	if namespace == "" {
		return nil
	}
	config, ok := c.apolloConfigCache.Load(namespace)
	if !ok {
		return nil
	}
	return config.(*Config)
}
