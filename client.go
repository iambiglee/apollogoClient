package agollo

import "github.com/apollogoClient/v1/env/config"

type Client interface {
	GetConfig(namespace string) *storage.Config
	GetConfigAndInit(namespace string) *storage.Config
	GetConfigCache(namespace string) agcache.CacheInterface
	GetDefaultConfigCache() agcache.CacheInterface
	GetApolloConfigCache() agcache.cacheInterface
	GetValue(key string) string
	GetStringValue(key string, defaultValue string) string
	GetIntValue(key string, defaultValue int) int
	GetFloatValue(key string, defaultValue float64) float64
	GetBoolValue(key string, defaultValue bool) bool
	GetStringSliceValue(key string, defaultValue []string) []string
	AddChangeListener(listener storage.ChangeListener)
	RemoveChangeListener(listener storage.ChangeListener)
}

func StartWithConfig(loadAppConfig func(*config.AppConfig, error)) (Client, error) {

}
