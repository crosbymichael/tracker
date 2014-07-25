package tracker

// Registry impements a persistent store for peers in the tracker
type Registry interface {
	// FetchPeers returns all the current peers in the tracker
	FetchPeers() ([]*Peer, error)

	// SavePeer saves the current peer in the registry with a specified ttl
	SavePeer(*Peer, int) error

	// DeletePeer removes the peer form the registry
	DeletePeer(*Peer) error
}
