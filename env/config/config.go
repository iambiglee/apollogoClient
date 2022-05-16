package config

import (
	"fmt"
	"github.com/apollogoClient/v1/utils"
	"net/url"
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

// ServerInfo Apollo 服务器信息
type ServerInfo struct {
	AppName     string `json:"appName"`
	InstanceID  string `json:"instanceID"`
	HomepageURL string `json:"homepageURL"`
	IsDone      bool   `json:"-"`
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

// GetServicesConfigURL 获取API的地址
// fmt.Sprintf是个啥子，理解为%s 替换为后面的参数
func (a *AppConfig) GetServicesConfigURL() string {
	return fmt.Sprintf("%sservices/config?appId=%s&ip=%s",
		a.GetHost(),
		url.QueryEscape(a.AppID),
		utils.GetInternal())
}

// GetHost 获取http 地址
// 为什么要这么做,这是strings的标准包
func (a *AppConfig) GetHost() string {
	u, err := url.Parse(a.IP)
	if err != nil {
		return a.IP
	}
	if !strings.HasSuffix(u.Path, "/") {
		return u.String() + "/"
	}
	return u.String()
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

// SetCurrentApolloConfig nolint
func (a *AppConfig) SetCurrentApolloConfig(apolloConfig *ApolloConnConfig) {
	a.currentConnApolloConfig.Set(apolloConfig.NamespaceName, apolloConfig)
}

type File interface {
	Load(fileName string, unmarshal func([]byte) (interface{}, error)) (interface{}, error)
	Write(content interface{}, configPath string) error
}
