package main

import (
	"log"
	"net/http"
	"time"

	"go_course/Homework_5/internal/api"
	"go_course/Homework_5/internal/config"
	"go_course/Homework_5/internal/storage"
)

// Command homework_15 — HTTP-сервис задач
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	var store api.Store
	if cfg.PGConn != "" {
		ps, err := storage.NewPostgresStore(cfg.PGConn)
		if err != nil {
			log.Printf("db not ready: %v, using memory", err)
			store = storage.NewMemoryStore()
		} else {
			store = ps
		}
	} else {
		store = storage.NewMemoryStore()
	}

	h := api.NewHandler(store)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listen %s", cfg.Addr)
	log.Fatal(srv.ListenAndServe())
}
