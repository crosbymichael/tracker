package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/crosbymichael/tracker/server"
)

var (
	interval    int
	minInterval int
	addr        string

	mux sync.Mutex
)

func init() {
	flag.StringVar(&addr, "addr", ":9090", "address of the tracker")
	flag.IntVar(&interval, "interval", 120, "interval for when Peers should poll for new peers")
	flag.IntVar(&minInterval, "min-interval", 30, "min poll interval for new peers")

	flag.Parse()
}

func main() {
	s := server.New(interval, minInterval, nil, logrus.New())

	if err := http.ListenAndServe(addr, s); err != nil {
		log.Fatal(err)
	}
}
