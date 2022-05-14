package remote

import (
	"github.com/apollogoClient/v1/env/config"
	"github.com/apollogoClient/v1/extension"
	"time"
)

const (
	//notify timeout
	notifyConnectTimeout = 10 * time.Minute //10m

	defaultContentKey = "content"
)

func loadBackupConfig(namespace string, appConfig config.AppConfig) []*config.ApolloConfig {
	appConfigs := make([]*config.ApolloConfig, 0)
	config.SplitNamespaces(namespace, func(namespace string) {
		c, err := extension.GetFileHandler().LoadConfigFile(appConfig.BackupConfigPath, appConfig.AppID, namespace)
		if err != nil || c == nil {
			return
		}
		appConfigs = append(appConfigs, c)
	})
	return appConfigs

}
