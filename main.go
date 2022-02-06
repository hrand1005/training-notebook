package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
	"github.com/hrand1005/training-notebook/handler"
)

func main() {
    logger := log.New(os.Stdout, "training-api", log.LstdFlags)
    setHandler := handler.NewSet(logger)

    serveMux := http.NewServeMux()
    serveMux.Handle("/sets", setHandler)

    server := &http.Server{
        Addr: ":8080",
        Handler: serveMux,
        IdleTimeout: 120*time.Second,
        ReadTimeout: 1*time.Second,
        WriteTimeout: 1*time.Second,
    }

    go func() {
        if err := server.ListenAndServe(); err != nil {
            logger.Fatal(err)
        }
    }()

    sigChan := make(chan os.Signal)
    signal.Notify(sigChan, os.Interrupt)
    signal.Notify(sigChan, os.Kill)

    sig := <-sigChan
    logger.Println("Recieved exit signal, executing graceful shutdown\nsignal: ", sig)

    timeout, _ := context.WithTimeout(context.Background(), 30*time.Second)
    server.Shutdown(timeout)
}
