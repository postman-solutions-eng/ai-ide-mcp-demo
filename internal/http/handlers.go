package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/ai-ide-mcp-demo/internal/models"
	"github.com/your-org/ai-ide-mcp-demo/internal/store"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	store     *store.Store
	startTime time.Time
}

// NewHandler creates a new Handler
func NewHandler(s *store.Store) *Handler {
	return &Handler{
		store:     s,
		startTime: time.Now(),
	}
}

// Health handles GET /health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Version:   "1.0.0",
		Uptime:    int64(time.Since(h.startTime).Seconds()),
	})
}

// ListServices handles GET /services
func (h *Handler) ListServices(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	owner := r.URL.Query().Get("owner")
	tier := r.URL.Query().Get("tier")
	runtime := r.URL.Query().Get("runtime")

	if limit <= 0 {
		limit = 50
	}

	services, total := h.store.ListServices(limit, owner, tier, runtime)

	json.NewEncoder(w).Encode(models.ServiceListResponse{
		Services: services,
		Pagination: models.Pagination{
			Limit: limit,
			Total: total,
		},
	})
}

// CreateService handles POST /services
func (h *Handler) CreateService(w http.ResponseWriter, r *http.Request) {
	var req models.CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON body",
		})
		return
	}

	if req.Name == "" || req.Owner == "" || req.Tier == "" || req.Runtime == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "validation_error",
			Message: "name, owner, tier, and runtime are required",
		})
		return
	}

	svc := h.store.CreateService(req)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(svc)
}

// GetService handles GET /services/{id}
func (h *Handler) GetService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	svc, err := h.store.GetService(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(svc)
}

// UpdateService handles PATCH /services/{id}
func (h *Handler) UpdateService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON body",
		})
		return
	}

	svc, err := h.store.UpdateService(id, req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(svc)
}

// GetEnvironment handles GET /services/{id}/environments/{env}
func (h *Handler) GetEnvironment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	env := chi.URLParam(r, "env")

	detail, err := h.store.GetEnvironment(id, env)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(detail)
}

// RotateKey handles POST /services/{id}/rotate-key
func (h *Handler) RotateKey(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Verify service exists
	_, err := h.store.GetService(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
		return
	}

	var req models.RotateKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON body",
		})
		return
	}

	now := time.Now()
	gracePeriod := time.Duration(req.GracePeriodMinutes) * time.Minute
	if gracePeriod == 0 {
		gracePeriod = 60 * time.Minute
	}

	json.NewEncoder(w).Encode(models.RotateKeyResponse{
		OperationID:       "op_" + uuid.New().String()[:8],
		ServiceID:         id,
		Environment:       req.Environment,
		KeyID:             "key_new_" + uuid.New().String()[:6],
		PreviousKeyID:     "key_old_" + uuid.New().String()[:6],
		GracePeriodEndsAt: now.Add(gracePeriod),
		ExpiresAt:         now.Add(90 * 24 * time.Hour),
		NotificationSent:  req.NotifyOwner,
		Status:            "in-progress",
	})
}

// GetEndpoints handles GET /services/{id}/endpoints
func (h *Handler) GetEndpoints(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	svc, err := h.store.GetService(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
		return
	}

	endpoints := make([]models.EndpointInfo, 0, len(svc.Environments))
	for envName, env := range svc.Environments {
		endpoints = append(endpoints, models.EndpointInfo{
			Environment: envName,
			BaseURL:     env.BaseURL,
			Region:      deriveRegion(envName),
			Auth:        normalizeAuth(env.AuthMode),
			Status:      deriveStatus(envName),
		})
	}

	json.NewEncoder(w).Encode(models.EndpointsResponse{
		ServiceID: id,
		Endpoints: endpoints,
	})
}

func deriveRegion(env string) string {
	switch env {
	case "prod":
		return "us-east-1"
	case "staging":
		return "us-west-2"
	default:
		return "us-west-2"
	}
}

func normalizeAuth(authMode string) string {
	switch authMode {
	case "api-key":
		return "apiKey"
	case "oauth2":
		return "oauth2"
	case "mtls":
		return "oauth2"
	case "":
		return "none"
	default:
		return authMode
	}
}

func deriveStatus(env string) string {
	if env == "prod" {
		return "healthy"
	}
	return "healthy"
}

