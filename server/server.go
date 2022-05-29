package server

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

// Server interface contains start method
type Server interface {
	Start(context.Context)
}

type server struct {
	h *http.Server
	c []io.Closer
}

// Start boots the server with the provided context.
func (s *server) Start(ctx context.Context) {
	// start server
	go func() {
		if err := s.h.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server encountered error while listening: %v\n", err)
		}
	}()

	log.Println("Server started")

	<-ctx.Done()

	log.Println("Shutting down server")

	// close outstanding server resources
	for _, v := range s.c {
		err := v.Close()
		if err != nil {
			log.Printf("Error closing server resource: %v", err)
		}
	}

	// shutdown gracefully with timeout context
	timeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.h.Shutdown(timeout); err != nil {
		log.Fatalf("Failed to shutdown correctly: %v\n", err)
	}
}
