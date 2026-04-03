package api

import (
	"net/http"
	"strconv"
)

// handleHealth is a simple liveness endpoint.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleChecks returns health check data.
//
// Without query params:  GET /api/checks → latest check for every URL
// With url param:        GET /api/checks?url=https://...&limit=50 → history for that URL
func (s *Server) handleChecks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// If a URL query param is provided, return history for that specific URL.
	urlParam := r.URL.Query().Get("url")
	if urlParam != "" {
		limit := 50
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil {
				limit = parsed
			}
		}

		checks, err := s.DB.GetCheckHistory(ctx, urlParam, limit)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to query check history")
			return
		}
		writeJSON(w, http.StatusOK, checks)
		return
	}

	// No URL param → return latest check for every monitored URL.
	checks, err := s.DB.GetLatestChecks(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query latest checks")
		return
	}
	writeJSON(w, http.StatusOK, checks)
}
