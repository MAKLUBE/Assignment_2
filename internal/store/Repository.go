package store

import "sync"

type Store[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewStore[K comparable, V any]() *Store[K, V] {
	return &Store[K, V]{
		data: make(map[K]V),
	}
}

func (s *Store[K, V]) Add(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

func (s *Store[K, V]) GetAll() []V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]V, 0, len(s.data))
	for _, v := range s.data {
		res = append(res, v)
	}
	return res
}
