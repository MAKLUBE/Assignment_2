package api

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/MAKLUBE/Assignment_2/internal/model"
	"github.com/MAKLUBE/Assignment_2/internal/store"
)

type Handler struct {
	Store *store.Store[string, *model.Task]
	Queue chan *model.Task
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
		w.Write([]byte("payload required"))
		return
	}

	task := &model.Task{
		Id:      nextID(),
		Payload: payload,
		Status:  model.Pending,
	}

	h.Store.Add(task.Id, task)
	h.Queue <- task

	w.Write([]byte("task created with id " + task.Id))
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	for _, task := range h.Store.GetAll() {
		w.Write([]byte(
			"ID=" + task.Id +
				" STATUS=" + string(task.Status) + "\n",
		))
	}
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/tasks/"):]

	task, ok := h.Store.Get(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte(
		"ID: " + task.Id +
			"\nPayload: " + task.Payload +
			"\nStatus: " + string(task.Status),
	))
}

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
