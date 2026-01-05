package api

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/MAKLUBE/Assignment_2/internal/model"
	"github.com/MAKLUBE/Assignment_2/internal/queue"
	"github.com/MAKLUBE/Assignment_2/internal/store"
)

type Handler struct {
	Store *store.Store
	Queue *queue.Queue[*model.Task]
}

var (
	idCounter int
	idMu      sync.Mutex
)

func nextID() int {
	idMu.Lock()
	defer idMu.Unlock()
	idCounter++
	return idCounter
}

// POST /tasks
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
	w.Write([]byte("task created with id " + strconv.Itoa(task.Id)))
}

// GET /tasks
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	for _, task := range h.Store.GetAll() {
		w.Write([]byte(
			"id=" + strconv.Itoa(task.Id) +
				" status=" + string(task.Status) + "\n",
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
		"ID: " + strconv.Itoa(task.Id) +
			"\nPayload: " + task.Payload +
			"\nStatus: " + string(task.Status),
	))
}

// GET /stats
func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	all := h.Store.GetAll()

	var submitted, inProgress, done int
	for _, t := range all {
		submitted++
		if t.Status == model.InProgress {
			inProgress++
		}
		if t.Status == model.Done {
			done++
		}
	}

	w.Write([]byte(
		"submitted: " + strconv.Itoa(submitted) + "\n" +
			"in_progress: " + strconv.Itoa(inProgress) + "\n" +
			"done: " + strconv.Itoa(done),
	))
}
