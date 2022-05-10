package env

import (
	"encoding/json"
	"github.com/apollogoClient/v1/env/config"
	jsonConfig "github.com/apollogoClient/v1/env/config/json"
	"github.com/apollogoClient/v1/utils"
	"os"
	"sync"
)

const (
	appConfigFile     = "app.properties"
	appConfigFilePath = "AGOLLO_CONF"

	defaultCluster   = "default"
	defaultNamespace = "application"
)

var executeConfigFileOnce sync.Once
var configFileExecutor config.File

func InitConfig(loadAppConfig func() (*config.AppConfig, error)) (*config.AppConfig, error) {
	return getLoadAppConfig(loadAppConfig)
}

func getLoadAppConfig(loadAppConfig func() (*config.AppConfig, error)) (*config.AppConfig, error) {
	if loadAppConfig != nil {
		return loadAppConfig()
	}
	configPath := os.Getenv(appConfigFilePath)
	if configPath == "" {
		configPath = appConfigFile
	}
	c, e := GetConfigFileExecutor().Load(configPath, Unmarshal)
	if c == nil {
		return nil, e
	}
	return c.(*config.AppConfig), e
}

// Unmarshal 反序列化
func Unmarshal(b []byte) (interface{}, error) {
	appConfig := &config.AppConfig{
		Cluster:        defaultCluster,
		NamespaceName:  defaultNamespace,
		IsBackupConfig: true,
	}
	err := json.Unmarshal(b, appConfig)
	if utils.IsNotNil(err) {
		return nil, err
	}
	appConfig.Init()
	return appConfig, nil
}

func GetConfigFileExecutor() config.File {
	executeConfigFileOnce.Do(func() {
		configFileExecutor = &jsonConfig.ConfigFile{}
	})
	return configFileExecutor
}
