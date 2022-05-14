package remote

import (
	"github.com/apollogoClient/v1/env"
	"github.com/apollogoClient/v1/env/config"
	"github.com/apollogoClient/v1/protocol/http"
	"strconv"
	"time"
)

// AbsApolloConfig 抽象Apollo配置
type AbsApolloConfig struct {
	remoteApollo ApolloConfig
}

// SyncWithNamespace 返回值怎么* 到底加在什么地方
func (a *AbsApolloConfig) SyncWithNamespace(namespace string, appConfigFunc func() config.AppConfig) *config.ApolloConfig {
	if appConfigFunc == nil {
		panic("cannot find apollo config,please check ")
	}
	appConfig := appConfigFunc()
	urlSuffix := a.remoteApollo.GetSyncURI(appConfig, namespace)

	c := &env.ConnectConfig{
		URI:     urlSuffix,
		AppID:   appConfig.AppID,
		Secrct:  appConfig.Secret,
		Timeout: notifyConnectTimeout,
	}
	if appConfig.SyncServerTimeout > 0 {
		duration, err := time.ParseDuration(strconv.Itoa(appConfig.SyncServerTimeout) + "s")
		if err != nil {
			return nil
		}
		c.Timeout = duration
	}
	callBack := a.remoteApollo.CallBack(namespace)
	apolloConfig, err := http.RequestRecovery(appConfig, c, &callBack)
	if err != nil || apolloConfig == nil {
		return nil
	}

	return apolloConfig.(*config.ApolloConfig)

}
