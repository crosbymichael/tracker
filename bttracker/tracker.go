package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/crosbymichael/tracker/registry"
	"github.com/crosbymichael/tracker/registry/inmem"
	"github.com/crosbymichael/tracker/registry/redis"
	"github.com/crosbymichael/tracker/server"
)

var (
	interval    int
	minInterval int
	addr        string
	redisAddr   string
	redisPass   string
	debug       bool

	mux sync.Mutex
)

func init() {
	flag.StringVar(&addr, "addr", ":80", "address of the tracker")
	flag.IntVar(&interval, "interval", 120, "interval for when Peers should poll for new peers")
	flag.IntVar(&minInterval, "min-interval", 30, "min poll interval for new peers")
	flag.BoolVar(&debug, "debug", false, "enable debug mode for logging")
	flag.StringVar(&redisAddr, "redis-addr", "", "address to a redis server for persistent peer data")
	flag.StringVar(&redisPass, "redis-pass", "", "password to use to connect to the redis server")

	flag.Parse()
}

func main() {
	var (
		logger   = logrus.New()
		registry registry.Registry
	)

	if debug {
		logger.Level = logrus.DebugLevel
	}

	if redisAddr != "" {
		registry = redis.New(redisAddr, redisPass)
	} else {
		registry = inmem.New()
	}

	s := server.New(interval, minInterval, registry, logger)
	if err := http.ListenAndServe(addr, s); err != nil {
		log.Fatal(err)
	}
}
