package env

import "time"

//ConnectConfig 网络请求配置
type ConnectConfig struct {
	Timeout time.Duration
	URI     string
	IsRetry bool
	AppID   string
	Secrct  string
}
