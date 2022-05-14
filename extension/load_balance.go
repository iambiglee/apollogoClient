package extension

import "github.com/apollogoClient/v1/cluster"

var defaultLoadBalance cluster.LoadBalance

func SetLoadBalance(loadBalance cluster.LoadBalance) {
	defaultLoadBalance = loadBalance
}

func GetLoadBalance() cluster.LoadBalance {
	return defaultLoadBalance
}
