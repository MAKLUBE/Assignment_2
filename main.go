package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Assignment2/internal/api"
	"Assignment2/internal/model"
	"Assignment2/internal/queue"
	"Assignment2/internal/store"
	"Assignment2/internal/worker"
)

func main() {

	taskStore := store.NewTaskStore()
	taskQueue := queue.NewQueue
	stopChan := make(chan struct{})

	handler := &api.Handler{
		Store: taskStore,
		Queue: taskQueue,
	}

	worker.StartWorker(1, taskQueue.Channel(), stopChan)
	worker.StartWorker(2, taskQueue.Channel(), stopChan)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s, ip, d := taskStore.Stats()
				log.Printf(
					"monitor: submitted=%d in_progress=%d done=%d",
					s, ip, d,
				)
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
		log.Println("server started on :8080")
		server.ListenAndServe()
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	close(stopChan)
	taskQueue.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	log.Println("server stopped gracefully")
}
