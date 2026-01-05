package model

import (
	"sync"
)

type Store struct {
	mu    sync.Mutex
	tasks map[string]*Task
}

func NewStore() *Store {
	return &Store{tasks: make(map[string]*Task)}
}

func (s *Store) Get(id string) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tasks[id]
}

func (s *Store) GetAll() []*Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	tasks := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}
