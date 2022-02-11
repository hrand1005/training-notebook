package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/handler"
)

func main() {

	router := gin.New()

	router.Use(logger())

	router.GET("/", handler.Home)

	router.GET("/sets", handler.ReadSets)
	router.GET("/sets/:id", handler.ReadSet)
	router.POST("/sets", handler.CreateSet)
	router.PUT("/sets/:id", handler.UpdateSet)
	router.DELETE("/sets/:id", handler.DeleteSet)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("Recieved exit signal, proceding with graceful shutdown:", sig)

	timeout, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeout)
}

func logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Processing %v request, URI: %v", c.Request.Method, c.Request.URL)
		t := time.Now()

		// before request
		c.Next()

		// after request
		latency := time.Since(t)
		log.Printf("Latency: %v\n", latency)

		status := c.Writer.Status()
		log.Printf("Response status: %v\n", status)
	}
}
