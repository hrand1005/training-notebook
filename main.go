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

	// create a group for the set resource
	setGroup := router.Group("/sets")

	// group routes that need the same set validation middleware
	setValidateGroup := setGroup.Group("")
	setValidateGroup.Use(handler.JSONSetValidator())

	setValidateGroup.POST("/", handler.CreateSet)
	setValidateGroup.PUT("/:id", handler.UpdateSet)

	setGroup.GET("", handler.ReadSets)
	setGroup.GET("/:id", handler.ReadSet)
	setGroup.DELETE("/:id", handler.DeleteSet)

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
