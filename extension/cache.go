package extension

import "github.com/apollogoClient/v1/agache"

var (
	globalCacheFactory agache.CacheFactory
)

// GetCacheFactory 获取CacheFactory
func GetCacheFactory() agache.CacheFactory {

	return globalCacheFactory
}

func SetCacheFactory(factory agache.CacheFactory) {
	globalCacheFactory = factory
}
