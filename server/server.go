package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/crosbymichael/tracker"
)

const bencodingFormat = "d8:intervali%de12:min intervali%de8:completei%de10:incompletei%de5:peersl%see"

type Server struct {
	interval    int
	minInterval int
	registry    tracker.Registry
	logger      *logrus.Logger

	mux *http.ServeMux
}

// New returns a new http.Handler for serving bittorrent tracker traffic
func New(interval, minInterval int, registry tracker.Registry, logger *logrus.Logger) http.Handler {
	s := &Server{
		interval:    interval,
		minInterval: minInterval,
		registry:    registry,
		logger:      logger,
		mux:         http.NewServeMux(),
	}

	s.mux.HandleFunc("/", s.tracker)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) tracker(w http.ResponseWriter, r *http.Request) {
	peer, err := tracker.PeerFromRequest(r)
	if err != nil {
		s.logger.WithField("error", err).Error("parsing peer from request")
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	switch event := r.URL.Query().Get("event"); event {
	case "stopped":
		s.logger.WithFields(logrus.Fields{
			"id":    peer.ID,
			"event": event,
		}).Info("received peer stop event")

		if err := s.registry.DeletePeer(peer); err != nil {
			s.logger.WithField("error", err).Error("remove peer from registry")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	default:
		s.logger.WithFields(logrus.Fields{
			"id":    peer.ID,
			"event": event,
		}).Info("received peer event")
	}

	if err := s.registry.SavePeer(peer, s.interval); err != nil {
		s.logger.WithField("error", err).Error("save peer from registry")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	peers, err := s.registry.FetchPeers()
	if err != nil {
		s.logger.WithField("error", err).Error("fetch peers from registry")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	var (
		completed int
		active    = []string{}
	)

	// build the bencoding strings for all the peers in the tracker
	for _, p := range peers {
		if p.IsSeed() {
			completed++

			// don't allow seeds to see each other
			if peer.IsSeed() {
				continue
			}
		}

		active = append(active, p.BTSerialize())
	}

	data := fmt.Sprintf(bencodingFormat, s.interval, s.minInterval, completed, len(active), strings.Join(active, ""))
	s.logger.WithField("data", data).Debug("response data")

	if _, err := fmt.Fprint(w, data); err != nil {
		s.logger.WithField("error", err).Error("write data to response")
	}
}
