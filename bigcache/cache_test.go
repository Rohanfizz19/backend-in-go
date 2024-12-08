package bigcache

import (
	"backend/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBigCache_GetValue(t *testing.T) {
	configs := config.BigCacheConfig{
		TTL:              5,
		Flag:             true,
		Shards:           16,
		MaxEntrySize:     20,
		StatsEnabled:     true,
		HardMaxCacheSize: 200,
	}
	cache, _ := NewInMemCache(configs)
	exists := cache.Exists("key1")
	assert.False(t, exists)

	err := cache.SetValue("key1", []byte("value1"))
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	val, found, _ := cache.GetValue("key1")
	assert.True(t, found)
	assert.Equal(t, []byte("value1"), val)

	time.Sleep(3 * time.Second)
	err = cache.SetValue("key2", []byte("value2"))
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	val, found, _ = cache.GetValue("key2")
	assert.True(t, found)
	assert.Equal(t, []byte("value2"), val)

	time.Sleep(4 * time.Second)
	val, found, _ = cache.GetValue("key1")
	assert.False(t, found)
	assert.Equal(t, []byte(nil), val)

	val, found, _ = cache.GetValue("key2")
	assert.True(t, found)
	assert.Equal(t, []byte("value2"), val)

	time.Sleep(3 * time.Second)
	val, found, _ = cache.GetValue("key2")
	assert.False(t, found)
	assert.Equal(t, []byte(nil), val)

}
