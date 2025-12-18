package store

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/postman-solutions-eng/ai-ide-mcp-demo/internal/models"
)

var (
	ErrNotFound = errors.New("service not found")
)

// Store provides in-memory storage for services
type Store struct {
	mu       sync.RWMutex
	services map[string]models.Service
	counter  int
}

// New creates a new in-memory store
func New() *Store {
	return &Store{
		services: make(map[string]models.Service),
	}
}

// ListServices returns all services with optional filtering
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
		result = append(result, svc)
		if len(result) >= limit {
			break
		}
	}

	return result, len(s.services)
}

// GetService retrieves a service by ID
func (s *Store) GetService(id string) (models.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	svc, ok := s.services[id]
	if !ok {
		return models.Service{}, ErrNotFound
	}
	return svc, nil
}

// CreateService adds a new service to the store
func (s *Store) CreateService(req models.CreateServiceRequest) models.Service {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	now := time.Now().UTC()

	svc := models.Service{
		ID:           fmt.Sprintf("svc_%06d", s.counter),
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

	s.services[svc.ID] = svc
	return svc
}

// UpdateService updates an existing service
func (s *Store) UpdateService(id string, req models.UpdateServiceRequest) (models.Service, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	svc, ok := s.services[id]
	if !ok {
		return models.Service{}, ErrNotFound
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
	svc.UpdatedAt = time.Now().UTC()

	s.services[id] = svc
	return svc, nil
}

// DeleteService removes a service from the store
func (s *Store) DeleteService(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.services[id]; !ok {
		return ErrNotFound
	}

	delete(s.services, id)
	return nil
}

