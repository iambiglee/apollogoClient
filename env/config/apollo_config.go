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
