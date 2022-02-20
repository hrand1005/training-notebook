package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/hrand1005/training-notebook/api/sets"
)

var prod = flag.Bool("prod", false, "Run production server")

func serve(prod bool, ctx context.Context) {

	// creates router with no default middleware, register logger
	router := gin.New()
	//router.Use(logger())

	// create a group for the set resource
	setGroup := router.Group("/sets")

	// the set handler contains CRUD operations for the set resource
	setHandler, err := sets.NewHandler()
	if err != nil {
		log.Fatal(err)
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

	// serve frontend build if production
	if prod {
		router.Use(static.Serve("/", static.LocalFile("./frontend/build", true)))
	}

	// configure server with gin router
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// start server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server encountered error while listening: %v\n", err)
		}
	}()

	log.Println("Server started")

	<-ctx.Done()

	log.Println("Shutting down server")

	// shutdown gracefully with timeout context
	timeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = server.Shutdown(timeout); err != nil {
		log.Fatalf("Failed to shutdown correctly: %v\n", err)
	}

	return
}

func main() {
	flag.Parse()

	// define shutdown conditions
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// define cancel context for server
	ctx, cancel := context.WithCancel(context.Background())

	// await kill signal
	go func() {
		sig := <-sigChan
		log.Printf("Recieved kill signal: %+v\n", sig)
		cancel()
	}()

	// start server with cancel context
	serve(*prod, ctx)
}

func logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		log.Printf("Processing %v request, URI: %v", c.Request.Method, c.Request.URL)
		t := time.Now()

		c.Next()

		// after request
		latency := time.Since(t)
		log.Printf("Latency: %v\n", latency)

		status := c.Writer.Status()
		log.Printf("Response status: %v\n", status)
	}
}
