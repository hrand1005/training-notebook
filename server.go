package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/hrand1005/training-notebook/api/sets"
	"gopkg.in/yaml.v3"
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
func buildServer(srvConf *serverConfig) (Server, error) {
	router := gin.New()

	// TODO: register resources in a different procedure
	// TODO: make resource registration configurable
	// create a group for the set resource
	setGroup := router.Group("/sets")

	// the set handler contains CRUD operations for the set resource
	setHandler, err := sets.New()
	if err != nil {
		return nil, err
	}

	// POST and PUT requests require JSONValidation
	setValidateGroup := setGroup.Group("")
	setValidateGroup.Use(setHandler.JSONValidator())
	setValidateGroup.POST("/", setHandler.Create)
	setValidateGroup.PUT("/:id", setHandler.Update)

	// GET and DELETE requests do not require JSONValidation
	setGroup.GET("", setHandler.ReadAll)
	setGroup.GET("/:id", setHandler.Read)
	setGroup.DELETE("/:id", setHandler.Delete)

	// go-openapi serve docs
	docOptions := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	// gin.WrapF converts http.HandlerFunc to gin HandlerFunc middleware
	docHandler := gin.WrapH(middleware.Redoc(docOptions, nil))
	router.GET("/docs", docHandler)
	router.StaticFile("/swagger.yaml", "./docs/swagger.yaml")

	if srvConf.Prod {
		registerStaticFiles(router)
	}

	s := &http.Server{
		Addr:         srvConf.Port,
		IdleTimeout:  srvConf.IdleTimeout,
		ReadTimeout:  srvConf.ReadTimeout,
		WriteTimeout: srvConf.WriteTimeout,
		Handler:      router,
	}

	return &server{
		s,
	}, nil
}

// registerStaticFiles should route our frontend builds
func registerStaticFiles(r *gin.Engine) {
	r.Use(static.Serve("/", static.LocalFile("./frontend/build", true)))
}

type serverConfig struct {
	Port         string        `yaml:"port"`
	IdleTimeout  time.Duration `yaml:"idle-timeout"`
	ReadTimeout  time.Duration `yaml:"read-timeout"`
	WriteTimeout time.Duration `yaml:"write-timeout"`
	Prod         bool
}

// loadServerConfig decodes a yaml configuration file, and sets deployment mode
func loadServerConfig(prodMode bool, configPath string) (*serverConfig, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// decode server config
	srvConf := &serverConfig{}
	d := yaml.NewDecoder(f)

	if err := d.Decode(srvConf); err != nil {
		return nil, err
	}

	srvConf.Prod = prodMode

	return srvConf, nil
}
