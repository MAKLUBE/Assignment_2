package main

import (
	"log"
	"net/http"

	"github.com/MAKLUBE/Assignment_2/internal/api"
	"github.com/MAKLUBE/Assignment_2/internal/model"
	"github.com/MAKLUBE/Assignment_2/internal/store"
	"github.com/MAKLUBE/Assignment_2/internal/worker"
)

func main() {
	taskStore := store.NewStore[string, *model.Task]()
	taskQueue := make(chan *model.Task, 10)

	handler := &api.Handler{
		Store: taskStore,
		Queue: taskQueue,
	}

	worker.StartWorker(1, taskQueue)
	worker.StartWorker(2, taskQueue)

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateTask(w, r)
		} else {
			handler.GetTasks(w, r)
		}
	})
	http.HandleFunc("/tasks/", handler.GetTask)
	http.HandleFunc("/stats", handler.Stats)

	log.Println("server started on :8080")
	http.ListenAndServe(":8080", nil)
}
