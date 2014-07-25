package tracker

import (
	"sync"
	"time"
)

// InMemRegistry implements a registry that stores Peer information in memory
// and is lost if the process is restarted.
type InMemRegistry struct {
	sync.Mutex

	peers map[string]*peerData
}

type peerData struct {
	p       *Peer
	expires time.Time
}

// NewInMemoryRegistry returns a new in memory registry for storing peer information
func NewInMemoryRegistry() Registry {
	return &InMemRegistry{
		peers: make(map[string]*peerData),
	}
}

func (r *InMemRegistry) FetchPeers() ([]*Peer, error) {
	r.Lock()

	var (
		out = []*Peer{}
		now = time.Now()
	)

	for _, p := range r.peers {
		if p.expires.After(now) {
			out = append(out, p.p)
		} else {
			key := r.getKey(p.p)
			delete(r.peers, key)
		}
	}

	r.Unlock()

	return out, nil
}

func (r *InMemRegistry) SavePeer(p *Peer, ttl int) error {
	r.Lock()

	key := r.getKey(p)
	r.peers[key] = &peerData{
		p:       p,
		expires: time.Now().Add(time.Duration(ttl) * time.Second),
	}

	r.Unlock()

	return nil
}

func (r *InMemRegistry) DeletePeer(p *Peer) error {
	r.Lock()

	key := r.getKey(p)
	delete(r.peers, key)

	r.Unlock()

	return nil
}

func (r *InMemRegistry) Close() error {
	return nil
}

func (r *InMemRegistry) getKey(p *Peer) string {
	return p.Hash()
}
