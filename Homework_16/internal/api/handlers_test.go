package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go_course/Homework_5/internal/model"

	"github.com/stretchr/testify/require"
)

type fakeStore struct {
	created []model.Task
	toGet   map[string][]model.Task
}

func (f *fakeStore) Create(task model.Task) error {
	f.created = append(f.created, task)
	return nil
}
func (f *fakeStore) GetByUID(uid string) ([]model.Task, error) {
	return f.toGet[uid], nil
}

func TestCreateTask_Unauthorized(t *testing.T) {
	h := NewHandler(&fakeStore{})
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{"title":"X","is_done":true}`))
	rr := httptest.NewRecorder()

	h.CreateTask(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestCreateTask_OK(t *testing.T) {
	fs := &fakeStore{}
	h := NewHandler(fs)

	body := bytes.NewBufferString(`{"title":"Do it","is_done":false}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	req.AddCookie(&http.Cookie{Name: "uid", Value: "u1"})
	rr := httptest.NewRecorder()

	h.CreateTask(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code)
	require.Len(t, fs.created, 1)
	require.Equal(t, "u1", fs.created[0].UID)
	require.Equal(t, "Do it", fs.created[0].Title)
	require.False(t, fs.created[0].IsDone)
}

func TestCreateTask_BadJSON(t *testing.T) {
	h := NewHandler(&fakeStore{})
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{bad-json`))
	req.AddCookie(&http.Cookie{Name: "uid", Value: "u1"})
	rr := httptest.NewRecorder()

	h.CreateTask(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetTasks_Unauthorized(t *testing.T) {
	h := NewHandler(&fakeStore{})
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	h.GetTasks(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetTasks_OK(t *testing.T) {
	fs := &fakeStore{
		toGet: map[string][]model.Task{
			"u1": {
				{ID: 1, UID: "u1", Title: "A", IsDone: false},
				{ID: 2, UID: "u1", Title: "B", IsDone: true},
			},
		},
	}
	h := NewHandler(fs)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	req.AddCookie(&http.Cookie{Name: "uid", Value: "u1"})
	rr := httptest.NewRecorder()

	h.GetTasks(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var got []model.Task
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	require.Len(t, got, 2)
	require.Equal(t, "A", got[0].Title)
	require.True(t, got[1].IsDone)
}
