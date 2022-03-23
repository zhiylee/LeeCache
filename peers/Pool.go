package peers

import (
	"LeeCache/consistentHash"
	"LeeCache/status"
	"sync"
)

// Pool : peers pool
type Pool struct {
	peers    map[string]*Peer
	hashRing *consistentHash.Map
	replicas int
	mux      sync.RWMutex
}

func NewPool(replicas int) *Pool {
	if replicas <= 0 {
		replicas = 50
	}

	return &Pool{
		peers:    make(map[string]*Peer),
		hashRing: consistentHash.New(replicas, nil),
		replicas: replicas,
	}
}

// Pick : pick the Peer where the cache located by key
func (p *Pool) Pick(key string) (*Peer, bool) {
	p.mux.RLock()
	peerAddr := p.hashRing.Get(key)
	p.mux.RUnlock()

	if peerAddr == "" {
		return &Peer{}, false
	}
	status.Log("Pick Peer %s", peerAddr)
	return p.peers[peerAddr], true
}

// peerIsExist : judge the Peer is exists by Peer.Addr
func (p *Pool) peerIsExist(addr string) bool {
	_, ok := p.peers[addr]
	return ok
}

// Add : add Peer
func (p *Pool) Add(peer *Peer) bool {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.peerIsExist(peer.Addr) {
		status.Log("add Peer %s failed,this Peer already exists", peer.Addr)
		return false
	}

	p.peers[peer.Addr] = peer
	p.hashRing.Add(peer.Addr)

	return true
}

// Delete : delete Peer
func (p *Pool) Delete(peer *Peer) bool {
	p.mux.Lock()
	defer p.mux.Unlock()

	if !p.peerIsExist(peer.Addr) {
		status.Log("delete Peer %s failed,this Peer is not exist")
		return false
	}

	delete(p.peers, peer.Addr)
	p.hashRing.Delete(peer.Addr)

	return true
}
