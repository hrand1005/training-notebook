package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/hrand1005/training-notebook/server"
	"github.com/joho/godotenv"
)

var prodMode = flag.Bool("prod", false, "Run in production mode and serve static files")
var configFile = flag.String("config", "", "Path to file containing server configs")

func main() {
	flag.Parse()
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env variables for signing key: %v", err)
	}

	// load server configs
	srvConf, err := loadConfig(*prodMode, *configFile)
	if err != nil {
		log.Fatalf("failed to load server configs from %v: %v", *configFile, err)
	}

	// build server with desired configuration
	server, err := ConstructHTTPServer(srvConf)
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

// ConstructHTTPServer directs the construction of the server using a serverBuilder
func ConstructHTTPServer(conf *Config) (server.Server, error) {
	serverBuilder := server.NewBuilder()

	serverBuilder.RegisterSwaggerDocs(conf.SwaggerSpec)
	serverBuilder.RegisterFileLogger(conf.LogFile)
	if conf.Prod {
		serverBuilder.RegisterFrontend(conf.Frontend)
	}
	serverBuilder.SetDB(conf.Database.Path)
	serverBuilder.SetServerAddr(conf.Server.Port)
	serverBuilder.SetIdleTimeout(conf.Server.IdleTimeout)
	serverBuilder.SetReadTimeout(conf.Server.ReadTimeout)
	serverBuilder.SetWriteTimeout(conf.Server.WriteTimeout)

	return serverBuilder.Construct()
}
