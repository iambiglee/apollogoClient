package serverlist

import (
	"encoding/json"
	"github.com/apollogoClient/v1/component"
	"github.com/apollogoClient/v1/env"
	"github.com/apollogoClient/v1/env/config"
	"github.com/apollogoClient/v1/env/server"
	"github.com/apollogoClient/v1/protocol/http"
	"strconv"
	"time"
)

// SyncServerIPListComponent Sync
type SyncServerIPListComponent struct {
	appConfig func() config.AppConfig
}

func (s SyncServerIPListComponent) Start() {
	SyncServerIPList(s.appConfig)
}

// SyncServerIPList 同步服务器列表的具体实现
func SyncServerIPList(appConfigFunc func() config.AppConfig) (map[string]*config.ServerInfo, error) {
	if appConfigFunc == nil {
		panic("can not find apollo config! please confirm")
	}
	appConfig := appConfigFunc()
	c := &env.ConnectConfig{
		AppID:  appConfig.AppID,
		Secrct: appConfig.Secret,
	}
	if appConfigFunc().SyncServerTimeout > 0 {
		duration, err := time.ParseDuration(strconv.Itoa(appConfigFunc().SyncServerTimeout))
		if err != nil {
			return nil, err
		}
		c.Timeout = duration
	}
	serverMap, err := http.Request(appConfig.GetServicesConfigURL(), c, &http.CallBack{
		SuccessCallBack: SyncServerIPListSuccessCallBack,
		AppConfigFunc:   appConfigFunc,
	})
	if serverMap == nil {
		return nil, err
	}

	m := serverMap.(map[string]*config.ServerInfo)
	server.SetServers(appConfig.GetHost(), m)
	return m, err
}

// SyncServerIPListSuccessCallBack 同步服务器成功之后的回调
//TODO json的Unmarshal用法
//TODO 这里要CallBack 做什么？
func SyncServerIPListSuccessCallBack(responseBody []byte, callback http.CallBack) (o interface{}, err error) {
	tmpServerInfo := make([]*config.ServerInfo, 0)
	err = json.Unmarshal(responseBody, &tmpServerInfo)

	if err != nil {
		return
	}

	if len(tmpServerInfo) == 0 {
		return
	}

	m := make(map[string]*config.ServerInfo)
	for _, serverInfo := range tmpServerInfo {
		if serverInfo == nil {
			continue
		}
		m[serverInfo.HomepageURL] = serverInfo
	}
	o = m
	return
}

// InitSyncServerIPList 初始化同步服务器信息列表
// TODO 但是下面用了& 是什么意思？
func InitSyncServerIPList(appConfig func() config.AppConfig) {
	go component.StartRefreshConfig(&SyncServerIPListComponent{appConfig})
}
