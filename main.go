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

	// creates router with no default middleware, register logger
	router := gin.New()
	router.Use(logger())

	// create a group for the set resource
	setGroup := router.Group("/sets")

	// the set handler contains CRUD operations for the set resource
	setHandler, err := handler.NewSet()
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
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// define shutdown conditions
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// handle shutdowns gracefully
	sig := <-sigChan
	log.Println("Recieved exit signal, proceding with graceful shutdown:", sig)

	timeout, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeout)
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
