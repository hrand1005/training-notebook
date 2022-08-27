package main

import (
	"context"
	"flag"
	"log"

	"github.com/joho/godotenv"

	"github.com/hrand1005/training-notebook/internal/app"
	"github.com/hrand1005/training-notebook/internal/config"
	"github.com/hrand1005/training-notebook/internal/http"
	"github.com/hrand1005/training-notebook/internal/http/handlers/user"
	"github.com/hrand1005/training-notebook/internal/mongo"
)

var configPath = flag.String("config", "", "Path to file containing server configs")

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file variables")
	}

	flag.Parse()
	conf, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load server configs from %s, err: %v", *configPath, err)
	}

	handle, err := mongo.New(conf.Database.URI, conf.Database.Name)
	if err != nil {
		log.Fatalf("failed to create new mongo db handle: %v", err)
	}
	defer handle.Close()

	userStore := mongo.NewUserStore(handle)
	userService := app.NewUserService(userStore)
	userHandler := user.NewUserHandler(userService, log.Default())

	server, err := http.BuildServer(conf.Server)
	if err != nil {
		log.Fatalf("failed to build server: %v", err)
	}

	server.RegisterHandler(userHandler)

	ctx := context.Background()
	server.Start(ctx)
}
