package extension

import "github.com/apollogoClient/v1/protocol/auth"

var authSign auth.HTTPAuth

//setHTTPAuth 设置Http验证
func SetHTTPAuth(http auth.HTTPAuth) {
	authSign = http
}

// GetHttPAuth 获取HttpAuth
func GetHttPAuth() auth.HTTPAuth {
	return authSign
}
