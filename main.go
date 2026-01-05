package main

import (
	"context"
	"github.com/MAKLUBE/Assignment_2/internal/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/tasks", Handler)
	http.HandleFunc("/stats", logsHandler)

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		log.Println("Starting server on port 8080")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start %w", err)
		}
	}()

	//graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	log.Println("Shutting down server")

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Server failed to shutdown %w", err)
	}

	log.Println("Server gracefully shutdown complete")

}
