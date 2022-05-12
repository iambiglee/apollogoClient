package server

import (
	"github.com/apollogoClient/v1/env/config"
	"sync"
)

var (
	ipMap      map[string]*Info
	serverLock sync.Mutex
	//next try connect period - 60 second
	nextTryConnectPeriod int64 = 30
)

type Info struct {
	//real servers ip
	serverMap       map[string]*config.ServerInfo
	nextTryConnTime int64
}

func SetServers(configIp string, serverMap map[string]*config.ServerInfo) {
	serverLock.Lock()
	defer serverLock.Unlock()
	ipMap[configIp] = &Info{
		serverMap: serverMap,
	}
}
