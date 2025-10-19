package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"go_course/Homework_5/internal/model"
)

// Store — контракт хранилища задач
type Store interface {
	Create(model.Task) (int64, error)
	GetByUID(uid string) ([]model.Task, error)
}

// Handler держит зависимости HTTP-хендлеров
type Handler struct {
	store Store
}

// NewHandler инициализирует обработчики
func NewHandler(s Store) *Handler {
	return &Handler{store: s}
}

// RegisterRoutes регистрирует endpoints
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tasks", h.GetTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)
}

// CreateTask — POST /tasks. Ждёт {"title": "...", "is_done": bool}
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	uid, err := readUID(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var in struct {
		Title  string `json:"title"`
		IsDone bool   `json:"is_done"`
	}
	if err := json.Unmarshal(raw, &in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}

	t := model.Task{UID: uid, Title: in.Title, IsDone: in.IsDone}
	id, err := h.store.Create(t)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	t.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(t)
}

// GetTasks — GET /tasks. Возвращает массив задач пользователя
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	uid, err := readUID(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	items, err := h.store.GetByUID(uid)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}

func readUID(r *http.Request) (string, error) {
	c, err := r.Cookie("uid")
	if err != nil {
		return "", err
	}
	v := strings.TrimSpace(c.Value)
	if v == "" {
		return "", errors.New("empty uid")
	}
	return v, nil
}
