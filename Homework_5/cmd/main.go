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

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateTask(w, r)
		} else if r.Method == http.MethodGet {
			handler.GetTasks(w, r)
		}
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
