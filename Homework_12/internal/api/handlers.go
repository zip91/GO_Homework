package api

import (
	"encoding/json"
	"net/http"

	"go_course/Homework_5/internal/model"
)

type Store interface {
	Create(task model.Task) error
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
	h.store.Create(t)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	uid := getUID(r)
	if uid == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	tasks, _ := h.store.GetByUID(uid)
	json.NewEncoder(w).Encode(tasks)
}

func getUID(r *http.Request) string {
	cookie, err := r.Cookie("uid")
	if err == nil {
		return cookie.Value
	}
	return ""
}
