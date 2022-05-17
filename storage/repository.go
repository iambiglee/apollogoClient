package storage

import (
	"container/list"
	"fmt"
	"github.com/apollogoClient/v1/agache"
	"github.com/apollogoClient/v1/env/config"
	"github.com/apollogoClient/v1/extension"
	"github.com/apollogoClient/v1/utils"
	"reflect"
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
//TODO 这里没写完
func (c *Cache) UpdateApolloConfig(apolloConfig *config.ApolloConfig, appConfigFunc func() config.AppConfig) {
	if apolloConfig == nil {
		return
	}
	appConfig := appConfigFunc()
	appConfig.SetCurrentApolloConfig(&apolloConfig.ApolloConnConfig)
	//cache := c.UpdateApolloConfigCache
	//cache(apolloConfig.Configurations, configCacheExpireTime, apolloConfig.NamespaceName)
	//appConfig.GetNotificationsMap().GetNotify(apolloConfig.NamespaceName)
	changeList := c.UpdateApolloConfigCache(apolloConfig.Configurations, configCacheExpireTime, apolloConfig.NamespaceName)
	notify := appConfig.GetNotificationsMap().GetNotify(apolloConfig.NamespaceName)
	c.pushNewestChanges(apolloConfig.NamespaceName, apolloConfig.Configurations, notify)
	if len(changeList) > 0 {
		event := createConfigChangeEvent(changeList, apolloConfig.NamespaceName, notify)
		c.pushChangeEvent(event)
	}
	if appConfig.GetIsBackupConfig() {
		//异步写入配置文件
		apolloConfig.AppID = appConfig.AppID
		go extension.GetFileHandler().WriteConfigFile(apolloConfig, appConfig.GetBackupConfigPath())
	}
}

//UpdateApolloConfigCache 根据conf[ig server返回的内容更新内存
func (c *Cache) UpdateApolloConfigCache(configurations map[string]interface{}, time int, namespace string) map[string]*ConfigChange {
	config := c.GetConfig(namespace)
	if config == nil {
		config = initConfig(namespace, extension.GetCacheFactory())
		c.apolloConfigCache.Store(namespace, config)
	}
	isInit := false

	defer func(c *Config) {
		if !isInit {
			return
		}
		b := c.GetIsInit()
		if b {
			return
		}
		c.isInit.Store(isInit)
		c.waitInit.Done()
	}(config)
	if (configurations == nil || len(configurations) == 0) && config.cache.EntryCount() == 0 {
		return nil
	}

	mp := map[string]bool{}
	config.cache.Range(func(key, value interface{}) bool {
		mp[key.(string)] = true
		return true
	})
	changes := make(map[string]*ConfigChange)
	if configurations != nil {
		for key, value := range configurations {
			if !mp[key] {
				changes[key] = createAddConfigChange(value)
			} else {
				//update
				oldValue, _ := config.cache.Get(key)
				if !reflect.DeepEqual(oldValue, value) {
					changes[key] = createModifyConfigChange(oldValue, value)
				}
			}
			err := config.cache.Set(key, value, time)
			if err != nil {
				fmt.Errorf("key:%s,value:%s is error", key, value)
			}
			delete(mp, key)
		}
	}
	for key := range mp {
		oldValue, _ := config.cache.Get(key)
		changes[key] = createDeletedConfigChange(oldValue)

		config.cache.Del(key)
	}
	isInit = true
	return changes

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

func (c *Cache) AddChangeListener(listener ChangeListener) {
	if listener == nil {
		return
	}
	c.ChangeListener.PushBack(listener)
}

//RemoveChangeListener 增加变更监控
func (c *Cache) RemoveChangeListener(listener ChangeListener) {
	if listener == nil {
		return
	}
	for i := c.ChangeListener.Front(); i != nil; i = i.Next() {
		apolloListener := i.Value.(ChangeListener)
		if listener == apolloListener {
			c.ChangeListener.Remove(i)
		}
	}
}

func (c *Cache) pushNewestChanges(name string, configurations map[string]interface{}, notify int64) {
	e := &FullChangeEvent{
		Changes: configurations,
	}
	e.Namespace = name
	e.NotificationID = notify
	c.pushChange(func(listener ChangeListener) {
		go listener.OnNewestChange(e)
	})

}

//TODO 为什么这里Func里面可以不写变量名，
func (c *Cache) pushChange(f func(ChangeListener)) {
	if c.ChangeListener == nil || c.ChangeListener.Len() == 0 {
		return
	}
	for i := c.ChangeListener.Front(); i != nil; i = i.Next() {
		listener := i.Value.(ChangeListener)
		f(listener)
	}
}

func (c *Cache) pushChangeEvent(event *ChangeEvent) {
	c.pushChange(func(listener ChangeListener) {
		go listener.OnChange(event)
	})
}

func (c *Config) GetCache() agache.CacheInterface {
	return c.cache
}

func GetDefaultNamespace() string {
	return defaultNamespace
}

func (c *Config) GetValue(key string) string {
	value := c.getConfigValue(key)
	if value == nil {
		return utils.Empty
	}

	v, ok := value.(string)
	if !ok {
		return utils.Empty
	}
	return v
}

//
func (c *Config) getConfigValue(key string) interface{} {
	b := c.GetIsInit()
	if !b {
		c.waitInit.Wait()
	}
	if c.cache == nil {
		return nil
	}
	value, err := c.cache.Get(key)
	if err != nil {
		return nil
	}
	return value

}

func (c *Config) GetIsInit() bool {
	return c.isInit.Load().(bool)
}

func (c *Config) GetStringValue(key string) string {
	value := c.GetValue(key)
	return value
}

//GetIntValue 获取配置值（int），获取不到则取默认值
func (c *Config) GetIntValue(key string, defaultValue int) int {
	value := c.getConfigValue(key)

	if value == nil {
		return defaultValue
	}
	v, ok := value.(int)
	if !ok {
		return defaultValue
	}
	return v
}

//GetFloatValue 获取配置值（float），获取不到则取默认值
func (c *Config) GetFloatValue(key string, defaultValue float64) float64 {
	value := c.getConfigValue(key)

	if value == nil {
		return defaultValue
	}

	v, ok := value.(float64)
	if !ok {
		return defaultValue
	}
	return v
}

//GetBoolValue 获取配置值（bool），获取不到则取默认值
func (c *Config) GetBoolValue(key string, defaultValue bool) bool {
	value := c.getConfigValue(key)
	v, ok := value.(bool)
	if !ok {
		return defaultValue
	}
	return v
}

func (c *Config) GetStringSliceValue(key string, defaultValue []string) []string {
	value := c.getConfigValue(key)
	v, ok := value.([]string)
	if !ok {
		return defaultValue
	}
	return v
}

func (c *Config) GetIntSliceValue(key string, defaultValue []int) []int {
	value := c.getConfigValue(key)
	v, ok := value.([]int)
	if !ok {
		return defaultValue
	}
	return v
}

//GetSliceValue 获取配置值（[]interface)
func (c *Config) GetSliceValue(key string, defaultValue []interface{}) []interface{} {
	value := c.getConfigValue(key)
	if value == nil {
		return defaultValue
	}
	v, ok := value.([]interface{})
	if !ok {
		return defaultValue
	}
	return v
}
