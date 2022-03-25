package peers

// Getter : get cache from Peer
type Getter interface {
	Get(group string, key string) ([]byte, error)
}

type Peer struct {
	Name   string
	Addr   string
	Status int
	Getter Getter
}

func (p *Peer) getStatus() string {
	return StatusText[p.Status]
}

func (p *Peer) isAvailable() bool {
	return p.Status == 200
}

func NewPeer(name string, addr string, a string) *Peer {
	var g Getter
	if a == "http" {
		g = &HTTPGetter{addr}
	} else {
		g = &RPCGetter{addr}
	}

	return &Peer{
		Name:   name,
		Addr:   addr,
		Status: StatusOK,
		Getter: g,
	}
}
