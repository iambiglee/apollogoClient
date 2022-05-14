package cluster

import "github.com/apollogoClient/v1/env/config"

type LoadBalance interface {
	Load(servers map[string]*config.ServerInfo) *config.ServerInfo
}
