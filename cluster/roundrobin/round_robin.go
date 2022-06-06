package roundrobin

import "github.com/apollogoClient/v1/env/config"

type RoundRobin struct {
}

// Load TODO 为什么之前不报错，为什么在我加了load之后才报错。这里实现的原理是什么
func (r *RoundRobin) Load(servers map[string]*config.ServerInfo) *config.ServerInfo {
	var returnServer *config.ServerInfo
	for _, server := range servers {
		if server.IsDone {
			continue
		}
		returnServer = server
		break
	}
	return returnServer
}
