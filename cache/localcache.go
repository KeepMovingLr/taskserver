package cache

import (
	"errors"
	"log"
	"sync"
)

var (
	MyCache = new(MyLocalCache)
	ErrSize = errors.New("cache size is wrong")
)

// define a local cache interface, different recycle value algorithm can implement this interface
type CacheInterface interface {
	AddToCache(key interface{}, value interface{})
	GetFromCache(key interface{}) (value interface{})
	InvalidCacheByKey(key interface{})
}

// this is an implement of CacheInterface by using LRU algorithm
type MyLocalCache struct {
	container LRUCache
	Lock      sync.Mutex
	MaxSize   int
}

// main method call this localCache to initialize MyLocalCache
func InitMyLocalCache(size int) (err error) {
	if size <= 0 {
		return ErrSize
	}
	MyCache.MaxSize = size
	MyCache.container = Constructor(MyCache.MaxSize)
	// initialize local cache success
	log.Println("init localcache end")
	return nil
}

func (c MyLocalCache) AddToCache(key, value interface{}) {
	c.Lock.Lock()
	c.container.Put(key, value)
	c.Lock.Unlock()
}

func (c MyLocalCache) GetFromCache(key interface{}) (value interface{}) {
	c.Lock.Lock()
	value = c.container.Get(key)
	c.Lock.Unlock()
	return value
}

func (c MyLocalCache) InvalidCacheByKey(key interface{}) {
	c.Lock.Lock()
	c.container.Invalid(key)
	c.Lock.Unlock()
}
