package tracker

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// NOTE: this is missing the peer id information
const peerFormat = "d2:ip%d:%s4:porti%dee"

// Peer represents a bittorrent peer
type Peer struct {
	ID        string `json:"id,omitempty"`
	IP        string `json:"ip,omitempty"`
	Port      int    `json:"port,omitempty"`
	InfoHash  string `json:"info_hash,omitempty"`
	Key       string `json:"key,omitempty"`
	BytesLeft int64  `json:"bytes_left,omitempty"`

	computedHash string
}

// IsSeed returns true if the peer has no more bytes left to receive
func (p *Peer) IsSeed() bool {
	return p.BytesLeft == 0
}

// BTSerialize returns the peer's information serialized in the the bencoding format
func (p *Peer) BTSerialize() string {
	return fmt.Sprintf(peerFormat, len(p.IP), p.IP, p.Port)
}

// PeerFromRequest returns a peer from an http GET request
func PeerFromRequest(r *http.Request) (*Peer, error) {
	v := r.URL.Query()

	port, err := strconv.Atoi(v.Get("port"))
	if err != nil {
		return nil, err
	}

	left, err := strconv.Atoi(v.Get("left"))
	if err != nil {
		return nil, err
	}

	p := &Peer{
		IP:        strings.Split(r.RemoteAddr, ":")[0], // we only need the ip not the port
		Port:      port,
		ID:        v.Get("peer_id"),
		InfoHash:  v.Get("info_hash"),
		Key:       v.Get("key"),
		BytesLeft: int64(left),
	}

	return p, nil
}

// Hash returns a sha1 of the peer ID and InfoHash
func (p *Peer) Hash() string {
	if p.computedHash == "" {
		hash := sha1.New()
		fmt.Fprint(hash, p.ID, p.InfoHash)

		p.computedHash = hex.EncodeToString(hash.Sum(nil))
	}

	return p.computedHash
}
