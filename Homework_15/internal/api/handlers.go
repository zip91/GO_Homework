package api

import (
	"encoding/json"
	"net/http"

	"go_course/Homework_5/internal/model"
)

type Store interface {
	Create(model.Task) error
	GetByUID(uid string) ([]model.Task, error)
}

type Handler struct {
	store Store
}

func NewHandler(s Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	uid := getUID(r)
	if uid == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var t model.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	t.UID = uid

	if err := h.store.Create(t); err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	uid := getUID(r)
	if uid == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tasks, err := h.store.GetByUID(uid)
	if err != nil {
		http.Error(w, "failed to get tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}

func getUID(r *http.Request) string {
	if cookie, err := r.Cookie("uid"); err == nil {
		return cookie.Value
	}
	return ""
}
