package registry

import (
	"github.com/crosbymichael/tracker/peer"
)

// Registry impements a persistent store for peers in the tracker
type Registry interface {
	// FetchPeers returns all the current peers in the tracker
	FetchPeers() ([]*peer.Peer, error)

	// SavePeer saves the current peer in the registry with a specified ttl
	SavePeer(*peer.Peer, int) error

	// DeletePeer removes the peer form the registry
	DeletePeer(*peer.Peer) error

	// Close closes any resources used by the registry
	Close() error
}
