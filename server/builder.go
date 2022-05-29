package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/hrand1005/training-notebook/api"
	"github.com/hrand1005/training-notebook/api/sets"
	"github.com/hrand1005/training-notebook/data"
)

// NewBuilder returns a builder with sensible server defaults.
func NewBuilder() *builder {
	return &builder{
		db: data.TestSetData,
		httpServer: http.Server{
			Addr: "8080",
		},
	}
}

// builder contains an interface for building a server.
// Call it's methods to configure a server, and call 'Construct' to instantiate it.
type builder struct {
	db              data.SetDB
	frontendPath    string
	logFile         string
	httpServer      http.Server
	swaggerSpecPath string
	err             error
}

func (b *builder) SetDB(dbPath string) {
	setDB, err := data.NewSetDB(dbPath)
	if err != nil {
		log.Printf("Encountered error building db, err: %v", err)
		b.err = err
		return
	}

	b.db = setDB
}

func (b *builder) RegisterSwaggerDocs(specPath string) {
	b.swaggerSpecPath = specPath
}

func (b *builder) RegisterFrontend(frontendPath string) {
	b.frontendPath = frontendPath
}

func (b *builder) RegisterFileLogger(logFile string) {
	b.logFile = logFile
}

func (b *builder) SetServerAddr(port string) {
	b.httpServer.Addr = port
}

func (b *builder) SetIdleTimeout(timeout time.Duration) {
	b.httpServer.IdleTimeout = timeout
}

func (b *builder) SetReadTimeout(timeout time.Duration) {
	b.httpServer.ReadTimeout = timeout
}

func (b *builder) SetWriteTimeout(timeout time.Duration) {
	b.httpServer.WriteTimeout = timeout
}

func (b *builder) Construct() (Server, error) {
	if b.err != nil {
		return nil, b.err
	}

	router := gin.New()

	// add outstanding resources to closers, and the server will call close
	// on them before gracefully shutting down
	closers := make([]io.Closer, 0, 5)

	// check if logger is set, if so use it with the router
	if b.logFile != "" {
		f, err := os.Create(b.logFile)
		if err != nil {
			log.Printf("Failed to create log file: %v\nProceeding without logger...", err)
		} else {
			logger := log.New(f, "", log.LstdFlags)
			router.Use(api.LatencyLogger(logger))
			closers = append(closers, f)
		}
	}

	setGroup := router.Group("/sets")
	// the set resource contains CRUD operations for sets
	// configure with SetDB, an interface for CRUD operations on set data
	setResource, err := sets.New(b.db)
	if err != nil {
		return nil, err
	}

	// registers resource CRUD operation handler funcs on the provided router group
	setResource.RegisterHandlers(setGroup)

	if b.swaggerSpecPath != "" {
		// set redoc options for swagger spec and create handler
		docOptions := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
		docHandler := gin.WrapH(middleware.Redoc(docOptions, nil))

		// create docs endpoint, register api documentation with SwaggerSpec
		docsGroup := router.Group("")
		docsGroup.GET("/docs", docHandler)
		docsGroup.StaticFile("/swagger.yaml", b.swaggerSpecPath)
	}

	if b.frontendPath != "" {
		router.Use(static.Serve("/", static.LocalFile(b.frontendPath, true)))
	}

	b.httpServer.Handler = router

	return &server{
		h: &b.httpServer,
		c: closers,
	}, nil
}

func (b *builder) Error() error {
	return b.err
}
