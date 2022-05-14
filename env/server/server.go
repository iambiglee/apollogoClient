package server

import (
	"github.com/apollogoClient/v1/env/config"
	"strings"
	"sync"
	"time"
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

func IsConnectDirectly(configIp string) bool {
	serverLock.Lock()
	defer serverLock.Unlock()
	s := ipMap[configIp]
	if s == nil || len(s.serverMap) == 0 {
		return false
	}
	if s.nextTryConnTime >= 0 && s.nextTryConnTime > time.Now().Unix() {
		return true
	}
	return false
}

func GetServers(configIp string) map[string]*config.ServerInfo {
	serverLock.Lock()
	defer serverLock.Unlock()
	if ipMap[configIp] == nil {
		return nil
	}
	return ipMap[configIp].serverMap
}

func SetDownNode(configIp string, host string) {
	serverLock.Lock()
	defer serverLock.Unlock()
	s := ipMap[configIp]
	if host == configIp {
		s.nextTryConnTime = nextTryConnectPeriod
	}
	for k, server := range s.serverMap {
		if strings.Index(k, host) > -1 {
			server.IsDone = true
		}
	}
}
