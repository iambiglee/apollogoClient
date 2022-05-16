package notify

import (
	"github.com/apollogoClient/v1/component/remote"
	"github.com/apollogoClient/v1/env/config"
	"github.com/apollogoClient/v1/storage"
	"time"
)

const (
	longPollInterval = 2 * time.Second //2s
)

type ConfigComponent struct {
	appConfigFunc func() config.AppConfig
	cache         *storage.Cache
}

func (c *ConfigComponent) SetAppConfig(appConfigFunc func() config.AppConfig) {
	c.appConfigFunc = appConfigFunc
}

func (c *ConfigComponent) SetCache(cache *storage.Cache) {
	c.cache = cache
}

func (c *ConfigComponent) Start() {
	timer := time.NewTimer(longPollInterval)
	instance := remote.CreateSyncApolloConfig()

	for {
		select {
		case <-timer.C:
			configs := instance.Sync(c.appConfigFunc)
			for _, apolloConfig := range configs {
				c.cache.UpdateApolloConfig(apolloConfig, c.appConfigFunc)
			}
			timer.Reset(longPollInterval)
		}
	}

}
