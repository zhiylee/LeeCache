package leecache

import (
	"LeeCache/peers"
	"LeeCache/status"
	"fmt"
	"log"
	"sync"
)

type Group struct {
	name   string
	getter Getter
	c      cache
	Peers  *peers.Pool
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
		status.Log("key = %s hit", key)
		return v, nil
	}

	return g.load(key)
}

// load : load value when not hit
func (g *Group) load(key string) (ByteView, error) {
	if g.Peers != nil {
		if peer, ok := g.Peers.Pick(key); ok && peer.Addr != Localhost {
			value, err := g.getFromPeer(peer, key)
			if err != nil {
				status.Log("failed to get cache from peer %s", peer.Addr)
			}

			return value, nil
		}
	}

	//get from localhost when no peer
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

// getFromPeer : get value from peer
func (g *Group) getFromPeer(peer *peers.Peer, key string) (ByteView, error) {
	bytes, err := peer.Getter.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}

	return ByteView{bytes}, nil
}

// GetGroup : get a group use name
func GetGroup(name string) *Group {
	mux.RLock()
	defer mux.RUnlock()

	return groups[name]
}
