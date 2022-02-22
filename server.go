package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/hrand1005/training-notebook/api/sets"
)

// Server interface contains start method
type Server interface {
	Start(context.Context)
}

type server struct {
	*http.Server
}

// Start boots the server with the provided context.
func (s *server) Start(ctx context.Context) {
	// start server
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server encountered error while listening: %v\n", err)
		}
	}()

	log.Println("Server started")

	<-ctx.Done()

	log.Println("Shutting down server")

	// shutdown gracefully with timeout context
	timeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(timeout); err != nil {
		log.Fatalf("Failed to shutdown correctly: %v\n", err)
	}
}

// buildServer configures server with desired configuration
func buildServer(conf *Config) (Server, error) {
	router := gin.New()

	// create a group for the set resource
	setGroup := router.Group("/sets")

	// the set resource contains CRUD operations for sets
	setResource, err := sets.New()
	if err != nil {
		return nil, err
	}

	// registers resource CRUD operation handler funcs on the provided router group
	setResource.RegisterHandlers(setGroup)

	// create docs endpoint, register api documentation with SwaggerSpec
	docsGroup := router.Group("")
	registerDocs(conf.SwaggerSpec, docsGroup)

	// we should serve the frontend if in production mode
	if conf.Prod {
		registerFrontend(router)
	}

	// create server using router and configs
	s := &http.Server{
		Addr:         conf.Server.Port,
		IdleTimeout:  conf.Server.IdleTimeout,
		ReadTimeout:  conf.Server.ReadTimeout,
		WriteTimeout: conf.Server.WriteTimeout,
		Handler:      router,
	}

	return &server{
		s,
	}, nil
}

// registerDocs creates documentation endpoints for our API using the provided
// swagger spec at specPath
func registerDocs(specPath string, g *gin.RouterGroup) {
	// go-openapi serve docs
	docOptions := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	// gin.WrapF converts http.HandlerFunc to gin HandlerFunc middleware
	docHandler := gin.WrapH(middleware.Redoc(docOptions, nil))
	g.GET("/docs", docHandler)
	g.StaticFile("/swagger.yaml", specPath)
}

// registerFrontend should route our frontend builds
func registerFrontend(r *gin.Engine) {
	r.Use(static.Serve("/", static.LocalFile("./frontend/build", true)))
}
