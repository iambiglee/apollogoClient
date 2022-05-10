package config

import (
	"strings"
	"sync"
)

var (
	defaultNotificationID = int64(-1)
	comma                 = ","
)

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
	notificationsMap        *notificationsMap
	currentConnApolloConfig *CurrentApolloConfig
}

// map[string]int64
type notificationsMap struct {
	notifications sync.Map
}

// Init 初始化 notificationsMap
func (a *AppConfig) Init() {
	a.currentConnApolloConfig = CreateCurrentApolloConfig()
	a.initAllNotifications(nil)
}

// InitAllNotifications 初始化notificationsMap
func (a *AppConfig) initAllNotifications(callback func(namespace string)) {
	ns := SplitNamespaces(a.NamespaceName, callback)
	a.notificationsMap = &notificationsMap{
		notifications: ns,
	}
}

//SplitNamespaces 根据namespace字符串分割后，并执行callback函数
func SplitNamespaces(name string, callback func(namespace string)) sync.Map {
	namespaces := sync.Map{}
	split := strings.Split(name, comma)
	for _, namespace := range split {
		if callback != nil {
			callback(namespace)
		}
		namespaces.Store(namespace, defaultNotificationID)
	}
	return namespaces

}

type File interface {
	Load(fileName string, unmarshal func([]byte) (interface{}, error)) (interface{}, error)
	Write(content interface{}, configPath string) error
}
