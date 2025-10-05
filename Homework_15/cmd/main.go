package main

import (
	"log"
	"net/http"

	"go_course/Homework_5/internal/api"
	"go_course/Homework_5/internal/config"
	"go_course/Homework_5/internal/storage"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var store api.Store
	psqlStore, err := storage.NewPostgresStore(cfg.PGConn)
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

	log.Printf("listening on %s", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, mux))
}
