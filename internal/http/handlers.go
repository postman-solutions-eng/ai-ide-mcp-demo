package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/postman-solutions-eng/ai-ide-mcp-demo/internal/models"
	"github.com/postman-solutions-eng/ai-ide-mcp-demo/internal/store"
)

var startTime = time.Now()

// Handler holds dependencies for HTTP handlers
type Handler struct {
	store *store.Store
}

// NewHandler creates a new Handler
func NewHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

// Health handles GET /health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "1.0.0",
		Uptime:    int64(time.Since(startTime).Seconds()),
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
	if services == nil {
		services = []models.Service{}
	}

	json.NewEncoder(w).Encode(models.ListServicesResponse{
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
			Message: "Service not found",
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
			Message: "Service not found",
		})
		return
	}

	json.NewEncoder(w).Encode(svc)
}

// DeleteService handles DELETE /services/{id}
func (h *Handler) DeleteService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.store.DeleteService(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "not_found",
			Message: "Service not found",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetEnvironment handles GET /services/{id}/environments/{env}
func (h *Handler) GetEnvironment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	envName := chi.URLParam(r, "env")

	svc, err := h.store.GetService(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "not_found",
			Message: "Service not found",
		})
		return
	}

	env, ok := svc.Environments[envName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "not_found",
			Message: "Environment not found",
		})
		return
	}

	json.NewEncoder(w).Encode(models.EnvironmentDetailsResponse{
		ServiceID:       svc.ID,
		ServiceName:     svc.Name,
		Environment:     envName,
		BaseURL:         env.BaseURL,
		AuthMode:        env.AuthMode,
		Config:          map[string]interface{}{"timeout": 30000, "retries": 3, "circuitBreaker": true},
		DeployedVersion: "v1.0.0",
		LastDeployedAt:  time.Now().UTC().Format(time.RFC3339),
	})
}

// RotateKey handles POST /services/{id}/rotate-key
func (h *Handler) RotateKey(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	svc, err := h.store.GetService(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "not_found",
			Message: "Service not found",
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

	if req.GracePeriodMinutes <= 0 {
		req.GracePeriodMinutes = 60
	}

	now := time.Now().UTC()
	gracePeriodEnd := now.Add(time.Duration(req.GracePeriodMinutes) * time.Minute)
	expiresAt := now.Add(90 * 24 * time.Hour) // 90 days

	json.NewEncoder(w).Encode(models.RotateKeyResponse{
		OperationID:       "op_" + uuid.New().String()[:8],
		ServiceID:         svc.ID,
		Environment:       req.Environment,
		KeyID:             "key_new_" + uuid.New().String()[:8],
		PreviousKeyID:     "key_old_" + uuid.New().String()[:8],
		GracePeriodEndsAt: gracePeriodEnd.Format(time.RFC3339),
		ExpiresAt:         expiresAt.Format(time.RFC3339),
		NotificationSent:  req.NotifyOwner,
		Status:            "in-progress",
	})
}

