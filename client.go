package agollo

import (
	"github.com/apollogoClient/v1/agache"
	"github.com/apollogoClient/v1/env"
	"github.com/apollogoClient/v1/env/config"
	storage "github.com/apollogoClient/v1/storatge"
)

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

func StartWithConfig(loadAppConfig func() (*config.AppConfig, error)) (Client, error) {
	//这里写了这么多，就是想找个合适的方法，获取到配置文件
	appConfig, err := env.InitConfig(loadAppConfig)

}
