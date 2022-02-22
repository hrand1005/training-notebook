package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

var prodMode = flag.Bool("prod", false, "Run in production mode and serve static files")
var configFile = flag.String("config", "", "Path to file containing server configs")

func main() {
	flag.Parse()

	// load server configs
	srvConf, err := loadConfig(*prodMode, *configFile)
	if err != nil {
		log.Fatalf("failed to load server configs from %v: %v", *configFile, err)
	}

	// build server with desired configuration
	server, err := buildServer(srvConf)
	if err != nil {
		log.Fatalf("failed to build server: %v", err)
	}

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
	server.Start(ctx)
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
