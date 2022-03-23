package leecache

import (
	"LeeCache/lru"
	"sync"
)

//封装lru.cache 添加互斥锁，实现并发安全

type cache struct {
	mux      sync.Mutex
	lru      *lru.Cache
	maxBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mux.Lock()
	defer c.mux.Unlock()

	//Lazy Initialization
	if c.lru == nil {
		c.lru = lru.New(c.maxBytes)
	}

	c.lru.Add(key, value)
}

func (c *cache) get(key string) (v ByteView, ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
