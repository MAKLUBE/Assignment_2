package store

import (
	"strconv"
	"sync"

	"github.com/MAKLUBE/Assignment_2/internal/model"
)

type Store struct {
	mu    sync.RWMutex
	tasks map[string]*model.Task
}

func NewStore() *Store {
	return &Store{
		tasks: make(map[string]*model.Task),
	}
}

func (s *Store) Add(task *model.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[strconv.Itoa(task.Id)] = task
}

func (s *Store) Get(id string) (*model.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	return task, ok
}

func (s *Store) GetAll() []*model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*model.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		res = append(res, task)
	}
	return res
}
