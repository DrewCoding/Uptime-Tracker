package api

import (
	"encoding/json"
	"log"
	"net/http"

	"tracker/internal/store"
)

// Server holds the HTTP server dependencies.
type Server struct {
	DB     *store.DB
	Router *http.ServeMux
}

// New creates a new API server with all routes registered.
func New(db *store.DB) *Server {
	s := &Server{
		DB:     db,
		Router: http.NewServeMux(),
	}
	s.routes()
	return s
}

// routes registers all API endpoints.
func (s *Server) routes() {
	s.Router.HandleFunc("GET /api/health", s.handleHealth)
	s.Router.HandleFunc("GET /api/checks", s.handleChecks)
}

// Handler returns the top-level HTTP handler with middleware applied.
func (s *Server) Handler() http.Handler {
	return cors(s.Router)
}

// cors is middleware that sets CORS headers to allow the React dev server
// (running on a different port) to call the API.
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// writeJSON marshals v to JSON and writes it to w.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
