package leecache

import (
	"fmt"
	"log"
	"sync"
)

type Group struct {
	name   string
	getter Getter
	c      cache
}

// mux : read and write mutex for var "groups"
var mux sync.RWMutex

// groups : a map save all group's pointer
var groups = make(map[string]*Group)

// NewGroup : create a new instance of Group
func NewGroup(name string, getter Getter, maxBytes int64) *Group {
	if name == "" || getter == nil {
		return nil
	}

	mux.Lock()
	defer mux.Unlock()

	g := &Group{
		name:   name,
		getter: getter,
		c:      cache{maxBytes: maxBytes},
	}
	groups[name] = g

	return g
}

// Get : get cache by key
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.c.get(key); ok {
		log.Printf("[LeeCache] key = %s hit\n", key)
		return v, nil
	}

	return g.load(key)
}

// load : load value when not hit
func (g *Group) load(key string) (ByteView, error) {
	return g.getFromLocally(key)
}

// getFromLocally : get value from localhost
func (g *Group) getFromLocally(key string) (ByteView, error) {
	log.Printf("[LeeCache] try load key = %s from Localhost", key)
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	v := ByteView{bytes}
	g.c.add(key, v)

	return v, nil
}

// GetGroup : get a group use name
func GetGroup(name string) *Group {
	mux.RLock()
	defer mux.RUnlock()

	return groups[name]
}
