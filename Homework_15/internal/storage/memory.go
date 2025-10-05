package storage

import (
	"sync"

	"go_course/Homework_5/internal/model"
)

type MemoryStore struct {
	mu   sync.RWMutex
	data []model.Task
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: []model.Task{}}
}

func (s *MemoryStore) Create(t model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, t)
	return nil
}

func (s *MemoryStore) GetByUID(uid string) ([]model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]model.Task, 0, len(s.data))
	for _, t := range s.data {
		if t.UID == uid {
			out = append(out, t)
		}
	}
	return out, nil
}
