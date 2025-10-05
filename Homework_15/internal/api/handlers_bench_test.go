package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go_course/Homework_5/internal/model"
)

type benchStore struct {
	data []model.Task
}

func (s *benchStore) Create(t model.Task) error {
	s.data = append(s.data, t)
	return nil
}
func (s *benchStore) GetByUID(uid string) ([]model.Task, error) {
	out := make([]model.Task, 0, len(s.data))
	for _, t := range s.data {
		if t.UID == uid {
			out = append(out, t)
		}
	}
	return out, nil
}

func BenchmarkCreateTask(b *testing.B) {
	s := &benchStore{}
	h := NewHandler(s)

	payload := map[string]any{
		"title":   "Do it",
		"is_done": false,
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
		req.AddCookie(&http.Cookie{Name: "uid", Value: "u1"})
		rr := httptest.NewRecorder()

		h.CreateTask(rr, req)

		if rr.Code != http.StatusCreated {
			b.Fatalf("unexpected status: %d", rr.Code)
		}
	}
}

func BenchmarkGetTasks(b *testing.B) {

	s := &benchStore{}
	for i := 0; i < 100; i++ {
		_ = s.Create(model.Task{UID: "u1", Title: "T", IsDone: i%2 == 0})
	}
	h := NewHandler(s)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		req.AddCookie(&http.Cookie{Name: "uid", Value: "u1"})
		rr := httptest.NewRecorder()

		h.GetTasks(rr, req)

		if rr.Code != http.StatusOK {
			b.Fatalf("unexpected status: %d", rr.Code)
		}

	}
}
