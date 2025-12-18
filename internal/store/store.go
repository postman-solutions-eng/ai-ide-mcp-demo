package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/ai-ide-mcp-demo/internal/models"
)

// Store provides in-memory storage for services
type Store struct {
	mu       sync.RWMutex
	services map[string]*models.Service
}

// New creates a new Store with optional seed data
func New() *Store {
	s := &Store{
		services: make(map[string]*models.Service),
	}
	s.seed()
	return s
}

func (s *Store) seed() {
	now := time.Now()
	s.services["svc_abc123"] = &models.Service{
		ID:      "svc_abc123",
		Name:    "payment-gateway",
		Owner:   "payments-team",
		Tier:    "gold",
		Runtime: "node",
		Tags:    []string{"pci-compliant", "critical"},
		Environments: map[string]models.Environment{
			"dev":  {BaseURL: "https://dev.payments.internal", AuthMode: "api-key"},
			"prod": {BaseURL: "https://payments.internal", AuthMode: "mtls"},
		},
		CreatedAt: now.Add(-5 * 24 * time.Hour),
		UpdatedAt: now.Add(-1 * time.Hour),
	}
	s.services["svc_def456"] = &models.Service{
		ID:      "svc_def456",
		Name:    "user-profile-service",
		Owner:   "identity-team",
		Tier:    "silver",
		Runtime: "python",
		Tags:    []string{"gdpr", "user-data"},
		Environments: map[string]models.Environment{
			"dev":     {BaseURL: "https://dev.users.internal", AuthMode: "api-key"},
			"staging": {BaseURL: "https://staging.users.internal", AuthMode: "oauth2"},
			"prod":    {BaseURL: "https://users.internal", AuthMode: "oauth2"},
		},
		CreatedAt: now.Add(-3 * 24 * time.Hour),
		UpdatedAt: now.Add(-12 * time.Hour),
	}
}

// ListServices returns services with optional filtering and pagination
func (s *Store) ListServices(limit int, owner, tier, runtime string) ([]models.Service, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 {
		limit = 50
	}

	var result []models.Service
	for _, svc := range s.services {
		if owner != "" && svc.Owner != owner {
			continue
		}
		if tier != "" && svc.Tier != tier {
			continue
		}
		if runtime != "" && svc.Runtime != runtime {
			continue
		}
		result = append(result, *svc)
	}

	total := len(result)
	if len(result) > limit {
		result = result[:limit]
	}

	return result, total
}

// GetService returns a service by ID
func (s *Store) GetService(id string) (*models.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	svc, ok := s.services[id]
	if !ok {
		return nil, fmt.Errorf("service not found: %s", id)
	}
	return svc, nil
}

// CreateService creates a new service
func (s *Store) CreateService(req models.CreateServiceRequest) *models.Service {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	id := "svc_" + uuid.New().String()[:8]

	svc := &models.Service{
		ID:           id,
		Name:         req.Name,
		Owner:        req.Owner,
		Tier:         req.Tier,
		Runtime:      req.Runtime,
		Description:  req.Description,
		Tags:         req.Tags,
		Environments: req.Environments,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	s.services[id] = svc
	return svc
}

// UpdateService updates an existing service
func (s *Store) UpdateService(id string, req models.UpdateServiceRequest) (*models.Service, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	svc, ok := s.services[id]
	if !ok {
		return nil, fmt.Errorf("service not found: %s", id)
	}

	if req.Tier != nil {
		svc.Tier = *req.Tier
	}
	if req.Description != nil {
		svc.Description = *req.Description
	}
	if req.Tags != nil {
		svc.Tags = req.Tags
	}
	svc.UpdatedAt = time.Now()

	return svc, nil
}

// GetEnvironment returns environment details for a service
func (s *Store) GetEnvironment(serviceID, envName string) (*models.EnvironmentDetailResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	svc, ok := s.services[serviceID]
	if !ok {
		return nil, fmt.Errorf("service not found: %s", serviceID)
	}

	env, ok := svc.Environments[envName]
	if !ok {
		return nil, fmt.Errorf("environment not found: %s", envName)
	}

	return &models.EnvironmentDetailResponse{
		ServiceID:       svc.ID,
		ServiceName:     svc.Name,
		Environment:     envName,
		BaseURL:         env.BaseURL,
		AuthMode:        env.AuthMode,
		Config:          models.EnvironmentConfig{Timeout: 30000, Retries: 3, CircuitBreaker: true},
		DeployedVersion: "v2.4.1",
		LastDeployedAt:  time.Now().Add(-24 * time.Hour),
	}, nil
}






