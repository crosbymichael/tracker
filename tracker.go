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

	clients = map[string]*client{}
	mux     sync.Mutex
)

func init() {
	flag.StringVar(&addr, "addr", ":9090", "address of the tracker")
	flag.IntVar(&interval, "interval", 120, "interval for when clients should poll for new peers")
	flag.IntVar(&minInterval, "min-interval", 30, "min poll interval for new peers")

	flag.Parse()
}

type client struct {
	addr     string
	peerID   string
	port     string
	ttl      int
	infoHash string
	key      string
	isSeed   bool
	event    string
}

func (c *client) id() string {
	return fmt.Sprintf("%s-%s", c.peerID, c.infoHash)
}

func (c *client) serialize() string {
	return fmt.Sprintf("d2:ip%d:%s4:porti%see", len(c.addr), c.addr, c.port)
}

func clientFromRequest(r *http.Request) *client {
	v := r.URL.Query()

	c := &client{
		addr:     strings.Split(r.RemoteAddr, ":")[0], // we only need the ip not the port
		port:     v.Get("port"),
		peerID:   v.Get("peer_id"),
		infoHash: v.Get("info_hash"),
		key:      v.Get("key"),
		isSeed:   v.Get("left") == "0",
		ttl:      interval,
		event:    v.Get("event"),
	}

	return c
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := clientFromRequest(r)

	if c.event == "stopped" {
		mux.Lock()
		delete(clients, c.id())
		mux.Unlock()

		return
	}

	mux.Lock()
	clients[c.id()] = c
	mux.Unlock()

	var (
		i, com int
		peers  = []string{}
	)

	mux.Lock()
	for _, cc := range clients {
		if cc.isSeed {
			com++

			if c.isSeed {
				continue
			}
		} else {
			i++
		}

		peers = append(peers, cc.serialize())
	}
	mux.Unlock()

	result := fmt.Sprintf("d8:intervali%de12:min intervali%de8:completei%de10:incompletei%de5:peersl%see",
		interval, minInterval, com, i, strings.Join(peers, ""))

	fmt.Fprint(w, result)
}

func main() {
	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
