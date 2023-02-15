package localcache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var _ LocalCache = (*cache.Cache)(nil)
var Cache LocalCache

type LocalCache interface {
	Set(k string, v interface{}, d time.Duration)
	Replace(k string, v interface{}, d time.Duration) error
	Get(k string) (interface{}, bool)
	Add(k string, v interface{}, d time.Duration) error
	Delete(k string)
}

func NewCache(expiration, interval time.Duration) LocalCache {
	return cache.New(expiration, interval)
}
