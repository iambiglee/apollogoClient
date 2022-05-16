package config

import "sync"

type CurrentApolloConfig struct {
	l       sync.RWMutex
	configs map[string]*ApolloConnConfig
}

func CreateCurrentApolloConfig() *CurrentApolloConfig {
	return &CurrentApolloConfig{
		configs: make(map[string]*ApolloConnConfig, 1),
	}
}

// ApolloConfig apollo配置
type ApolloConfig struct {
	ApolloConnConfig
	Configurations map[string]interface{} `json:"configurations"`
}

func (c *CurrentApolloConfig) Set(namespace string, config *ApolloConnConfig) {
	c.l.Lock()
	defer c.l.Unlock()
	c.configs[namespace] = config
}
