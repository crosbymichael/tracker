package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/crosbymichael/tracker"
	"github.com/crosbymichael/tracker/server"
)

var (
	interval    int
	minInterval int
	addr        string
	debug       bool

	mux sync.Mutex
)

func init() {
	flag.StringVar(&addr, "addr", ":9090", "address of the tracker")
	flag.IntVar(&interval, "interval", 120, "interval for when Peers should poll for new peers")
	flag.IntVar(&minInterval, "min-interval", 30, "min poll interval for new peers")
	flag.BoolVar(&debug, "debug", false, "enable debug mode for logging")

	flag.Parse()
}

func main() {
	logger := logrus.New()
	if debug {
		logger.Level = logrus.Debug
	}

	s := server.New(interval, minInterval, tracker.NewInMemoryRegistry(), logger)

	if err := http.ListenAndServe(addr, s); err != nil {
		log.Fatal(err)
	}
}
