package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MAKLUBE/Assignment_2/internal/api"
	"github.com/MAKLUBE/Assignment_2/internal/queue"
	"github.com/MAKLUBE/Assignment_2/internal/store"
	"github.com/MAKLUBE/Assignment_2/internal/worker"
)

func main() {
	taskStore := store.NewStore()
	stopChan := make(chan struct{})
	taskQueue := queue.NewQueue
	
	handler := &api.Handler{
		Store: taskStore,
		Queue: taskQueue,
	}

	worker.StartWorker(1, taskQueue.Channel(), stopChan)
	worker.StartWorker(2, taskQueue.Channel(), stopChan)

	// monitor
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("monitor: tasks =", len(taskStore.GetAll()))
			case <-stopChan:
				return
			}
		}
	}()

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateTask(w, r)
		} else {
			handler.GetTasks(w, r)
		}
	})

	http.HandleFunc("/tasks/", handler.GetTask)
	http.HandleFunc("/stats", handler.Stats)

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		log.Println("Server started on :8080")
		server.ListenAndServe()
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	close(stopChan)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	log.Println("Server stopped gracefully")
}
