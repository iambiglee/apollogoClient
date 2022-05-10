package config

import "sync"

type notificationMap struct {
	notifications sync.Map
}
type CurrentApolloConfig struct {
	l       sync.RWMutex
	configs map[string]*ApolloConnConfig
}
type ApolloConnConfig struct {
	AppID         string `json:"appID"`
	Cluster       string `json:"cluster"`
	NamespaceName string `json:"namespaceName"`
	ReleaseKey    string `json:"releaseKey"`
	sync.RWMutex
}

type AppConfig struct {
	AppID                   string `json:"appID"`
	Cluster                 string `json:"cluster"`
	NamespaceName           string `json:"namespaceName"`
	IP                      string `json:"IP"`
	IsBackupConfig          bool   `json:"isBackupConfig"`
	BackupConfigPath        string `json:"backupConfigPath"`
	Secret                  string `json:"secret"`
	Label                   string `json:"label"`
	SyncServerTimeout       int    `json:"syncServerTimeout"`
	MustStart               bool   `json:"mustStart"`
	notificationsMap        *notificationMap
	currentConnApolloConfig *CurrentApolloConfig
}
