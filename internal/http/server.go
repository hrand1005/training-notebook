package http

import (
	"context"
	"io"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hrand1005/training-notebook/internal/config"
	"github.com/hrand1005/training-notebook/internal/http/handlers"
)

// Server interfaces can be started with a context.
// The context should be used to send the server a 'done' signal.
type Server interface {
	Start(context.Context)
	RegisterHandler(handlers.Handler)
}

type server struct {
	api       fiber.Router
	app       *fiber.App
	addr      string
	resources []io.Closer
}

// BuildServer creates a Server with the provided configs.
func BuildServer(conf config.ServerConfig) (Server, error) {
	app := fiber.New(
		fiber.Config{
			ReadTimeout:  conf.ReadTimeout,
			WriteTimeout: conf.WriteTimeout,
			IdleTimeout:  conf.IdleTimeout,
		},
	)

	return &server{
		api:  app.Group(conf.API.Prefix),
		app:  app,
		addr: conf.Port,
	}, nil
}

// RegisterHandler registers the given handler on the server's api endpoint.
func (s *server) RegisterHandler(h handlers.Handler) {
	h.Register(s.api)
}

// Start starts the server listening on the configured address. The
// server runs until it recieves a done signal on the provided context,
// at which point it executes graceful shutdown.
func (s *server) Start(ctx context.Context) {
	go func() {
		if err := s.app.Listen(s.addr); err != nil {
			log.Fatalf("server encountered error: %v", err)
		}
	}()

	log.Printf("server started")

	<-ctx.Done()

	log.Printf("server recieved done signal, executing graceful shutdown")

	for _, v := range s.resources {
		if err := v.Close(); err != nil {
			log.Printf("error closing server resource: %v", err)
		}
	}

	if err := s.app.Shutdown(); err != nil {
		log.Fatalf("server failed to shutdown correctly: %v", err)
	}
}
