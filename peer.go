package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Peer struct {
	addr     string
	peerID   string
	port     string
	ttl      int
	infoHash string
	key      string
	isSeed   bool
	event    string
}

func (p *Peer) id() string {
	return fmt.Sprintf("%s-%s", p.peerID, p.infoHash)
}

func (p *Peer) serialize() string {
	return fmt.Sprintf("d2:ip%d:%s4:porti%see", len(p.addr), p.addr, p.port)
}

func PeerFromRequest(r *http.Request) *Peer {
	v := r.URL.Query()

	p := &Peer{
		addr:     strings.Split(r.RemoteAddr, ":")[0], // we only need the ip not the port
		port:     v.Get("port"),
		peerID:   v.Get("peer_id"),
		infoHash: v.Get("info_hash"),
		key:      v.Get("key"),
		isSeed:   v.Get("left") == "0",
		ttl:      interval,
		event:    v.Get("event"),
	}

	return p
}
