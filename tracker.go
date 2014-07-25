package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	interval    int
	minInterval int
	addr        string

	peers = map[string]*Peer{}
	mux   sync.Mutex
)

func init() {
	flag.StringVar(&addr, "addr", ":9090", "address of the tracker")
	flag.IntVar(&interval, "interval", 120, "interval for when Peers should poll for new peers")
	flag.IntVar(&minInterval, "min-interval", 30, "min poll interval for new peers")

	flag.Parse()
}

func handler(w http.ResponseWriter, r *http.Request) {
	peer := PeerFromRequest(r)

	if peer.event == "stopped" {
		mux.Lock()
		delete(peers, peer.id())
		mux.Unlock()

		return
	}

	mux.Lock()
	peers[peer.id()] = peer
	mux.Unlock()

	var (
		i, com int
		alive  = []string{}
	)

	mux.Lock()
	for _, cc := range peers {
		if cc.isSeed {
			com++

			if peer.isSeed {
				continue
			}
		} else {
			i++
		}

		alive = append(alive, cc.serialize())
	}
	mux.Unlock()

	result := fmt.Sprintf("d8:intervali%de12:min intervali%de8:completei%de10:incompletei%de5:peersl%see",
		interval, minInterval, com, i, strings.Join(alive, ""))

	fmt.Fprint(w, result)
}

func main() {
	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
