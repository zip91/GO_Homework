package storage

import (
	"sync"
	"sync/atomic"

	"go_course/Homework_5/internal/model"
)

// MemoryStore — in-memory реализация хранилища
type MemoryStore struct {
	mu    sync.RWMutex
	seq   int64
	byUID map[string][]model.Task
}

// NewMemoryStore возвращает пустое хранилище в памяти
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{byUID: make(map[string][]model.Task)}
}

// Create добавляет запись и присваивает id
func (m *MemoryStore) Create(t model.Task) (int64, error) {
	id := atomic.AddInt64(&m.seq, 1)
	t.ID = id

	m.mu.Lock()
	m.byUID[t.UID] = append(m.byUID[t.UID], t)
	m.mu.Unlock()

	return id, nil
}

// GetByUID отдаёт копию среза задач пользователя
func (m *MemoryStore) GetByUID(uid string) ([]model.Task, error) {
	m.mu.RLock()
	src := m.byUID[uid]
	out := make([]model.Task, len(src))
	copy(out, src)
	m.mu.RUnlock()
	return out, nil
}
