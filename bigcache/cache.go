package bigcache

import (
	"backend/config"
	"backend/metrics"
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

type Cache interface {
	GetValue(key string) ([]byte, bool, error)
	SetValue(key string, value []byte) error
	Exists(key string) bool
}

type BigCache struct {
	cache   *bigcache.BigCache
	metrics metrics.BigCacheMetrics
}

func NewInMemCache(configs config.BigCacheConfig) (*BigCache, error) {

	cacheMetrics := metrics.InitBigCacheMetrics()

	config := bigcache.DefaultConfig(time.Duration(configs.TTL) * time.Second)
	config.Shards = configs.Shards
	config.StatsEnabled = configs.StatsEnabled
	config.MaxEntrySize = configs.MaxEntrySize         // in bytes
	config.HardMaxCacheSize = configs.HardMaxCacheSize // in MB

	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		return nil, err
	}

	go metrics.CollectStatistics(cache, cacheMetrics)

	return &BigCache{
		cache:   cache,
		metrics: cacheMetrics,
	}, nil
}

func (c *BigCache) GetValue(key string) (resp []byte, found bool, err error) {
	resp, err = c.cache.Get(key)
	if err != nil {
		// Check if the error is ErrEntryNotFound
		if err == bigcache.ErrEntryNotFound {
			return nil, false, nil
		} else {
			return nil, false, err
		}
	}
	return resp, true, nil
}

func (c *BigCache) SetValue(key string, value []byte) error {
	return c.cache.Set(key, value)
}

func (c *BigCache) Exists(key string) bool {
	_, found, _ := c.GetValue(key)
	return found
}
