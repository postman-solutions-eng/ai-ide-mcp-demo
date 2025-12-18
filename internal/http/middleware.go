package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/your-org/ai-ide-mcp-demo/internal/models"
)

const (
	HeaderRequestID = "X-Platform-Request-ID"
	HeaderTraceID   = "X-Platform-Trace-ID"
)

// RequireRequestID middleware validates X-Platform-Request-ID header
func RequireRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(HeaderRequestID)
		if requestID == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:   "missing_header",
				Message: "X-Platform-Request-ID header is required",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AddTraceHeaders middleware echoes request ID and adds trace ID to responses
func AddTraceHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(HeaderRequestID)
		if requestID != "" {
			w.Header().Set(HeaderRequestID, requestID)
		}
		w.Header().Set(HeaderTraceID, uuid.New().String())
		next.ServeHTTP(w, r)
	})
}

// JSONContentType middleware sets Content-Type to application/json
func JSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}






