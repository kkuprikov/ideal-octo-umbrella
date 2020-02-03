// Package streamprocessor routes
package streamprocessor

import (
	"context"
	"log"
	"net/http"
	"sync"
)

// Start method for our server
func (s *Server) Start(ctx context.Context, port string, wg *sync.WaitGroup, tasks *SafeTasks) {
	s.Router.POST("/api/stream_stats/start", s.statsStart(ctx, wg, tasks))
	s.Router.POST("/api/stream_stats/stop", s.statsStop(ctx, wg, tasks))
	s.Router.POST("/api/snapshots/start", s.screenshotsStart(ctx, wg, tasks))
	s.Router.POST("/api/snapshots/stop", s.screenshotsStop(ctx, wg, tasks))

	if err := http.ListenAndServe(":"+port, s.Router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
