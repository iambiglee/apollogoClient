package json

import "github.com/apollogoClient/v1/env/config"

type FileHandler struct {
}

func (f *FileHandler) WriteConfigFile(config *config.ApolloConfig, configPath string) error {
	//TODO implement me
	panic("implement me")
}

func (f *FileHandler) GetConfigFile(configDir string, appID string, namespace string) string {
	//TODO implement me
	panic("implement me")
}

func (f *FileHandler) LoadConfigFile(configDir string, appID string, namespace string) (*config.ApolloConfig, error) {
	//TODO implement me
	panic("implement me")
}
