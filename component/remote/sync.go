package remote

import (
	"encoding/json"
	"fmt"
	"github.com/apollogoClient/v1/constant"
	"github.com/apollogoClient/v1/env/config"
	"github.com/apollogoClient/v1/extension"
	"github.com/apollogoClient/v1/protocol/http"
	"github.com/apollogoClient/v1/utils"
	"net/url"
	"path"
)

// TODO 不用申明名称吗
type syncApolloConfig struct {
	AbsApolloConfig
}

func (s syncApolloConfig) GetNotifyURLSuffix(notifications string, config config.AppConfig) string {
	return ""
}

func (s syncApolloConfig) GetSyncURI(config config.AppConfig, namespaceName string) string {
	return fmt.Sprintf("configfiles/json/%s/%s/%s?&ip=%s&label=%s",
		url.QueryEscape(config.AppID),
		url.QueryEscape(config.Cluster),
		url.QueryEscape(namespaceName),
		utils.GetInternal(),
		url.QueryEscape(config.Label))
}

func (s syncApolloConfig) CallBack(namespace string) http.CallBack {
	return http.CallBack{
		SuccessCallBack:   processJSONFile,
		NotModifyCallBack: touchApolloConfigCache,
		Namespace:         namespace,
	}
}

func touchApolloConfigCache() error {
	return nil
}

// processJSONFile TODO 返回值可以命名，也可以不命名？60行左右Parse 没有返回值也行？Parse是个接口，没有实现类，为什么不报错
//解析Json, 将key 值返回
func processJSONFile(bytes []byte, back http.CallBack) (o interface{}, err error) {
	apolloConfig := &config.ApolloConfig{}
	apolloConfig.NamespaceName = back.Namespace

	configurations := make(map[string]interface{}, 0)
	apolloConfig.Configurations = configurations
	err = json.Unmarshal(bytes, &apolloConfig.Configurations)
	if err != nil {
		return nil, err
	}
	parser := extension.GetFormatParser(constant.ConfigFileFormat(path.Ext(apolloConfig.NamespaceName)))
	if parser == nil {
		parser = extension.GetFormatParser(constant.DEFAULT)
	}
	if parser == nil {
		return apolloConfig, nil
	}
	parse, err := parser.Parse(configurations[defaultContentKey])
	if err != nil {
		fmt.Errorf("error")
	}
	if len(parse) > 0 {
		apolloConfig.Configurations = parse
	}
	return apolloConfig, nil
}

// Sync 配置信息的优先级：Apollo-->Local File.都没有就报错
func (a syncApolloConfig) Sync(appConfigFunc func() config.AppConfig) []*config.ApolloConfig {
	appConfig := appConfigFunc()
	configs := make([]*config.ApolloConfig, 0, 8)
	config.SplitNamespaces(appConfig.NamespaceName, func(namespace string) {
		apolloConfig := a.SyncWithNamespace(namespace, appConfigFunc)
		if apolloConfig != nil {
			configs = append(configs, apolloConfig)
			return
		}
		configs = append(configs, loadBackupConfig(appConfig.NamespaceName, appConfig)...)
	})
	return configs
}

// CreateSyncApolloConfig 创建同步获取 Apollo 配置
//TODO 这里报错为什么要说没有实现类
func CreateSyncApolloConfig() ApolloConfig {
	a := &syncApolloConfig{}
	a.remoteApollo = a
	return a
}
