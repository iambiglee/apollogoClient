package remote

import (
	"github.com/apollogoClient/v1/env/config"
	"github.com/apollogoClient/v1/protocol/http"
)

// ApolloConfig Apollo 配置 TODO 怎么这里又来一个
type ApolloConfig interface {
	// GetNotifyURLSuffix 获取异步更新路劲
	GetNotifyURLSuffix(notifications string, config config.AppConfig) string

	// GetSyncURI 获取同步路径
	GetSyncURI(config config.AppConfig, namespaceName string) string

	//Sync 同步获取Apollo 配置
	Sync(appConfigFunc func() config.AppConfig) []*config.AppConfig

	// CallBack 根据namespace获取callback数据
	CallBack(namespace string) http.CallBack

	// SyncWithNamespace 通过namespace 同步Apollo
	SyncWithNamespace(namespace string, appConfigFunc func() config.AppConfig) *config.AppConfig
}
