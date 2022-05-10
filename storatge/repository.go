package storage

import (
	"github.com/apollogoClient/v1/agache"
	"sync"
	"sync/atomic"
)

type Config struct {
	namespace string
	cache     agache.CacheInterface
	isInit    atomic.Value
	waitInit  sync.WaitGroup
}
