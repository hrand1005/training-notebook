package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/hrand1005/training-notebook/handler"
)

func main() {

	router := mux.NewRouter()

	// register middleware on router
	router.Use(logger)

	// register handlers on router
	router.HandleFunc("/sets", handler.ReadSets).Methods(http.MethodGet)
	router.HandleFunc("/sets/{id:[0-9]+}", handler.ReadSet).Methods(http.MethodGet)
	router.HandleFunc("/sets", handler.CreateSet).Methods(http.MethodPost)
	router.HandleFunc("/sets/{id:[0-9]+}", handler.UpdateSet).Methods(http.MethodPut)
	router.HandleFunc("/sets/{id:[0-9]+}", handler.DeleteSet).Methods(http.MethodDelete)

	// configure server
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

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("Processing %v request, URI: %v\n", r.Method, r.RequestURI)
		t := time.Now()
		// assuming no more middleware...
		next.ServeHTTP(rw, r)

		latency := time.Since(t)
		log.Printf("Latency: %v\n", latency)
	})
}
