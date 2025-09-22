package main

import (
	"log"
	"net/http"
	"os"

	"go_course/Homework_5/internal/api"
	"go_course/Homework_5/internal/storage"
)

func main() {
	connStr := os.Getenv("PG_CONN")

	var store api.Store
	psqlStore, err := storage.NewPostgresStore(connStr)
	if err != nil {
		log.Println("PostgreSQL unavailable, fallback to memory")
		store = storage.NewMemoryStore()
	} else {
		store = psqlStore
	}

	handler := api.NewHandler(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateTask(w, r)
		case http.MethodGet:
			handler.GetTasks(w, r)
		default:
			w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
