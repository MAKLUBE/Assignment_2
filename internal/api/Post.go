package api

import (
	"net/http"
	"strconv"
	"sync"

	"Assignment2/internal/model"
	"Assignment2/internal/queue"
	"Assignment2/internal/store"
)

type Handler struct {
	Store *store.Store
	Queue *queue.Queue[*model.Task]
}

var (
	idCounter int
	idMu      sync.Mutex
)

func nextID() string {
	idMu.Lock()
	defer idMu.Unlock()
	idCounter++
	return strconv.Itoa(idCounter)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	payload := r.FormValue("payload")
	if payload == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("payload is required"))
		return
	}

	task := &model.Task{
		Id:      nextID(),
		Payload: payload,
		Status:  model.Pending,
	}

	h.Store.Add(task)
	h.Queue.Push(task)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("task created with id " + task.ID))
}

// GET /tasks
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	for _, task := range h.Store.GetAll() {
		w.Write([]byte(
			"ID: " + task.ID +
				" | Status: " + string(task.Status) + "\n",
		))
	}
}

// GET /tasks/{id}
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/tasks/"):]

	task, ok := h.Store.Get(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("task not found"))
		return
	}

	w.Write([]byte(
		"ID: " + task.ID +
			"\nPayload: " + task.Payload +
			"\nStatus: " + string(task.Status),
	))
}

// GET /stats
func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	sub, prog, done := h.Store.Stats()

	w.Write([]byte(
		"submitted: " + strconv.Itoa(sub) + "\n" +
			"in_progress: " + strconv.Itoa(prog) + "\n" +
			"completed: " + strconv.Itoa(done),
	))
}
