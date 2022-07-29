package server

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/hrand1005/training-notebook/api"
	"github.com/hrand1005/training-notebook/data"

	_ "github.com/mattn/go-sqlite3"
)

// NewBuilder returns a builder with sensible server defaults.
func NewBuilder() *builder {
	return &builder{
		httpServer: http.Server{
			Addr: "8080",
		},
	}
}

// builder contains an interface for building a server.
// Call it's methods to configure a server, and call 'Construct' to instantiate it.
type builder struct {
	db              *sql.DB
	frontendPath    string
	logFile         string
	httpServer      http.Server
	swaggerSpecPath string
	err             error
}

func (b *builder) SetDB(dbPath string) {
	db, err := data.SqliteDB(dbPath)
	if err != nil {
		b.err = err
		return
	}

	b.db = db
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
	apiGroup := router.Group("/api")

	// add outstanding resources to closers, and the server will call close
	// on them before gracefully shutting down
	closers := []io.Closer{b.db}

	// check if logger is set, if so use it with the router
	if b.logFile != "" {
		f, err := os.Create(b.logFile)
		if err != nil {
			log.Printf("Failed to create log file: %v\nProceeding without logger...", err)
		} else {
			logger := log.New(f, "", log.LstdFlags)
			apiGroup.Use(api.LatencyLogger(logger))
			closers = append(closers, f)
		}
	}

	if b.db == nil {
		// TODO: initialize test mode db
		return nil, fmt.Errorf("no db found, test mode not yet implemented")
	}

	if err := api.RegisterAll(b.db, apiGroup); err != nil {
		return nil, fmt.Errorf("registering api endpoints: %v", err)
	}

	if b.swaggerSpecPath != "" {
		// set redoc options for swagger spec and create handler
		docOptions := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
		docHandler := gin.WrapH(middleware.Redoc(docOptions, nil))

		// create docs endpoint, register api documentation with SwaggerSpec
		apiGroup.GET("/docs", docHandler)
		apiGroup.StaticFile("/swagger.yaml", b.swaggerSpecPath)
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
