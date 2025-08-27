package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

var (
	store = map[string]*Todo{}
	mu    sync.RWMutex
	idSeq = 0
)

func listTodos(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]*Todo, 0, len(store))
	for _, t := range store {
		out = append(out, t)
	}
	writeJSON(w, http.StatusOK, out)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title string `json:"title"`
		Done  bool   `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(in.Title) == "" {
		http.Error(w, "title required", http.StatusUnprocessableEntity)
		return
	}
	idSeq++
	t := &Todo{
		ID:        fmtID(idSeq),
		Title:     in.Title,
		Done:      in.Done,
		CreatedAt: time.Now().UTC(),
	}
	mu.Lock()
	store[t.ID] = t
	mu.Unlock()
	writeJSON(w, http.StatusCreated, t)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/todos/")
	mu.RLock()
	t, ok := store[id]
	mu.RUnlock()
	if !ok {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func fmtID(n int) string {
	return "todo-" + strconv.Itoa(n)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func gzipRequestDecompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Header.Get("Content-Encoding")) == "gzip" && r.Body != nil {
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "invalid gzip body", http.StatusBadRequest)
				return
			}
			r.Body = struct {
				io.Reader
				io.Closer
			}{Reader: gr, Closer: closerChain{gr, r.Body}}
			r.Header.Del("Content-Length")
		}
		next.ServeHTTP(w, r)
	})
}

type closerChain []io.Closer

func (cc closerChain) Close() error {
	for _, c := range cc {
		_ = c.Close()
	}
	return nil
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gz      *gzip.Writer
	enabled bool
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.enabled {
		return w.ResponseWriter.Write(b)
	}
	return w.gz.Write(b)
}

func (w *gzipResponseWriter) Close() {
	if w.gz != nil {
		w.gz.Close()
	}
}

func gzipResponseCompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		ww := &gzipResponseWriter{ResponseWriter: w}
		ct := w.Header().Get("Content-Type")
		if strings.HasPrefix(ct, "application/json") || strings.HasPrefix(ct, "text/html") {
			gz := gzip.NewWriter(w)
			ww.gz = gz
			ww.enabled = true
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Add("Vary", "Accept-Encoding")
			defer ww.Close()
		}

		next.ServeHTTP(ww, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listTodos(w, r)
		case http.MethodPost:
			createTodo(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/todos/", getTodo)

	handler := gzipRequestDecompressor(gzipResponseCompressor(mux))

	log.Println("listening :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
