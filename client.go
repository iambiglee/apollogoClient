package agollo

import (
	"errors"
	"github.com/apollogoClient/v1/agache"
	"github.com/apollogoClient/v1/component"
	"github.com/apollogoClient/v1/component/notify"
	"github.com/apollogoClient/v1/component/remote"
	"github.com/apollogoClient/v1/component/serverlist"
	"github.com/apollogoClient/v1/env"
	"github.com/apollogoClient/v1/env/config"
	storage "github.com/apollogoClient/v1/storage"
)

var syncApolloConfig = remote.CreateSyncApolloConfig()

type Client interface {
	GetConfig(namespace string) *storage.Config
	GetConfigAndInit(namespace string) *storage.Config
	GetConfigCache(namespace string) agache.CacheInterface
	GetDefaultConfigCache() agache.CacheInterface
	GetApolloConfigCache() agache.CacheInterface
	GetValue(key string) string
	GetStringValue(key string, defaultValue string) string
	GetIntValue(key string, defaultValue int) int
	GetFloatValue(key string, defaultValue float64) float64
	GetBoolValue(key string, defaultValue bool) bool
	GetStringSliceValue(key string, defaultValue []string) []string
	AddChangeListener(listener storage.ChangeListener)
	RemoveChangeListener(listener storage.ChangeListener)
}

// internalClient apollo 客户端实例
type internalClient struct {
	initAppConfigFunc func() (*config.AppConfig, error)
	appConfig         *config.AppConfig
	cache             *storage.Cache
}

func (c *internalClient) GetConfig(namespace string) *storage.Config {
	return c.GetConfigAndInit(namespace)
}

//GetConfigAndInit 根据namespace获取apollo配置
func (c *internalClient) GetConfigAndInit(namespace string) *storage.Config {
	if namespace == "" {
		return nil
	}

	config := c.cache.GetConfig(namespace)

	if config == nil {
		//init cache
		storage.CreateNamespaceConfig(namespace)

		//sync config
		syncApolloConfig.SyncWithNamespace(namespace, c.getAppConfig)
	}

	config = c.cache.GetConfig(namespace)

	return config
}

//GetConfigCache 根据namespace获取apollo配置的缓存
func (c *internalClient) GetConfigCache(namespace string) agache.CacheInterface {
	config := c.GetConfigAndInit(namespace)
	if config == nil {
		return nil
	}
	return config.GetCache()
}

func (c *internalClient) GetDefaultConfigCache() agache.CacheInterface {
	configAndInit := c.GetConfigAndInit(storage.GetDefaultNamespace())
	if configAndInit != nil {
		return configAndInit.GetCache()
	}
	return nil
}

func (c *internalClient) GetApolloConfigCache() agache.CacheInterface {
	return c.GetDefaultConfigCache()
}

func (c *internalClient) GetValue(key string) string {
	return c.GetConfig(storage.GetDefaultNamespace()).GetValue(key)
}

func (c *internalClient) GetStringValue(key string, defaultValue string) string {
	return c.GetConfig(storage.GetDefaultNamespace()).GetStringValue(key)
}

func (c *internalClient) GetIntValue(key string, defaultValue int) int {
	return c.GetConfig(storage.GetDefaultNamespace()).GetIntValue(key, defaultValue)
}

func (c *internalClient) GetFloatValue(key string, defaultValue float64) float64 {
	return c.GetConfig(storage.GetDefaultNamespace()).GetFloatValue(key, defaultValue)
}

func (c *internalClient) GetBoolValue(key string, defaultValue bool) bool {
	return c.GetConfig(storage.GetDefaultNamespace()).GetBoolValue(key, defaultValue)
}

func (c *internalClient) GetStringSliceValue(key string, defaultValue []string) []string {
	return c.GetConfig(storage.GetDefaultNamespace()).GetStringSliceValue(key, defaultValue)
}

func (c *internalClient) AddChangeListener(listener storage.ChangeListener) {
	c.cache.AddChangeListener(listener)
}

func (c *internalClient) RemoveChangeListener(listener storage.ChangeListener) {
	c.cache.RemoveChangeListener(listener)
}

func (c *internalClient) getAppConfig() config.AppConfig {
	return *c.appConfig
}

func StartWithConfig(loadAppConfig func() (*config.AppConfig, error)) (Client, error) {
	//这里写了这么多，就是想找个合适的方法，获取到配置文件
	appConfig, err := env.InitConfig(loadAppConfig)
	if err != nil {
		return nil, err
	}
	//创造出client
	c := create()
	if appConfig != nil {
		c.appConfig = appConfig
	}
	c.cache = storage.CreateNamespaceConfig(appConfig.NamespaceName)
	appConfig.Init()
	serverlist.InitSyncServerIPList(c.getAppConfig)

	//fist sync
	configs := syncApolloConfig.Sync(c.getAppConfig)
	if len(configs) == 0 && appConfig != nil && appConfig.MustStart {
		return nil, errors.New("no config")
	}
	for _, apolloConfig := range configs {
		c.cache.UpdateApolloConfig(apolloConfig, c.getAppConfig)
	}

	//开始长轮训
	configComponent := &notify.ConfigComponent{}
	configComponent.SetAppConfig(c.getAppConfig)
	configComponent.SetCache(c.cache)
	go component.StartRefreshConfig(configComponent)

	return c, err
}

//create
//这里不清楚为什么要用&，以及括号里面为什么是两个一样的东西
//应为，如果返回值是一个*， 那么return 就要是内存地址
func create() *internalClient {
	appConfig := env.InitFileConfig()
	return &internalClient{
		appConfig: appConfig,
	}
}
