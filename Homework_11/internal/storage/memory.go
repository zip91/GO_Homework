package storage

import "go_course/Homework_5/internal/model"

type MemoryStore struct {
	data []model.Task
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: []model.Task{}}
}

func (s *MemoryStore) Create(task model.Task) error {
	s.data = append(s.data, task)
	return nil
}

func (s *MemoryStore) GetByUID(uid string) ([]model.Task, error) {
	var result []model.Task
	for _, t := range s.data {
		if t.UID == uid {
			result = append(result, t)
		}
	}
	return result, nil
}
