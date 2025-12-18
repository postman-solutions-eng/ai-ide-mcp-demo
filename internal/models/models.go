package models

import "time"

// Service represents a platform service in the catalog
type Service struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Owner        string                 `json:"owner"`
	Tier         string                 `json:"tier"`
	Runtime      string                 `json:"runtime"`
	Description  string                 `json:"description,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	Environments map[string]Environment `json:"environments,omitempty"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

// Environment represents environment-specific configuration
type Environment struct {
	BaseURL  string `json:"baseUrl"`
	AuthMode string `json:"authMode"`
}

// ServiceListResponse is the response for GET /services
type ServiceListResponse struct {
	Services   []Service  `json:"services"`
	Pagination Pagination `json:"pagination"`
}

// Pagination metadata for list responses
type Pagination struct {
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor,omitempty"`
	Total  int    `json:"total"`
}

// CreateServiceRequest is the request body for POST /services
type CreateServiceRequest struct {
	Name         string                 `json:"name"`
	Owner        string                 `json:"owner"`
	Tier         string                 `json:"tier"`
	Runtime      string                 `json:"runtime"`
	Description  string                 `json:"description,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	Environments map[string]Environment `json:"environments,omitempty"`
}

// UpdateServiceRequest is the request body for PATCH /services/{id}
type UpdateServiceRequest struct {
	Tier        *string  `json:"tier,omitempty"`
	Description *string  `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// EnvironmentDetailResponse is the response for GET /services/{id}/environments/{env}
type EnvironmentDetailResponse struct {
	ServiceID       string            `json:"serviceId"`
	ServiceName     string            `json:"serviceName"`
	Environment     string            `json:"environment"`
	BaseURL         string            `json:"baseUrl"`
	AuthMode        string            `json:"authMode"`
	Config          EnvironmentConfig `json:"config"`
	DeployedVersion string            `json:"deployedVersion"`
	LastDeployedAt  time.Time         `json:"lastDeployedAt"`
}

// EnvironmentConfig holds environment-specific settings
type EnvironmentConfig struct {
	Timeout        int  `json:"timeout"`
	Retries        int  `json:"retries"`
	CircuitBreaker bool `json:"circuitBreaker"`
}

// RotateKeyRequest is the request body for POST /services/{id}/rotate-key
type RotateKeyRequest struct {
	Environment        string `json:"environment"`
	NotifyOwner        bool   `json:"notifyOwner"`
	GracePeriodMinutes int    `json:"gracePeriodMinutes"`
}

// RotateKeyResponse is the response for POST /services/{id}/rotate-key
type RotateKeyResponse struct {
	OperationID       string    `json:"operationId"`
	ServiceID         string    `json:"serviceId"`
	Environment       string    `json:"environment"`
	KeyID             string    `json:"keyId"`
	PreviousKeyID     string    `json:"previousKeyId"`
	GracePeriodEndsAt time.Time `json:"gracePeriodEndsAt"`
	ExpiresAt         time.Time `json:"expiresAt"`
	NotificationSent  bool      `json:"notificationSent"`
	Status            string    `json:"status"`
}

// HealthResponse is the response for GET /health
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    int64     `json:"uptime"`
}

// ErrorResponse is a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// EndpointsResponse is the response for GET /services/{id}/endpoints
type EndpointsResponse struct {
	ServiceID string         `json:"serviceId"`
	Endpoints []EndpointInfo `json:"endpoints"`
}

// EndpointInfo represents a single endpoint for a service environment
type EndpointInfo struct {
	Environment string `json:"environment"`
	BaseURL     string `json:"baseUrl"`
	Region      string `json:"region"`
	Auth        string `json:"auth"`
	Status      string `json:"status"`
}

