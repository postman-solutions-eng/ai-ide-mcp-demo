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

// Environment represents environment-specific config for a service
type Environment struct {
	BaseURL  string `json:"baseUrl"`
	AuthMode string `json:"authMode"`
}

// HealthResponse is returned by GET /health
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	Uptime    int64  `json:"uptime"`
}

// ListServicesResponse is returned by GET /services
type ListServicesResponse struct {
	Services   []Service  `json:"services"`
	Pagination Pagination `json:"pagination"`
}

// Pagination metadata for list responses
type Pagination struct {
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor,omitempty"`
	Total  int    `json:"total"`
}

// CreateServiceRequest is the body for POST /services
type CreateServiceRequest struct {
	Name         string                 `json:"name"`
	Owner        string                 `json:"owner"`
	Tier         string                 `json:"tier"`
	Runtime      string                 `json:"runtime"`
	Description  string                 `json:"description,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	Environments map[string]Environment `json:"environments,omitempty"`
}

// UpdateServiceRequest is the body for PATCH /services/{id}
type UpdateServiceRequest struct {
	Tier        *string  `json:"tier,omitempty"`
	Description *string  `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// EnvironmentDetailsResponse is returned by GET /services/{id}/environments/{env}
type EnvironmentDetailsResponse struct {
	ServiceID       string                 `json:"serviceId"`
	ServiceName     string                 `json:"serviceName"`
	Environment     string                 `json:"environment"`
	BaseURL         string                 `json:"baseUrl"`
	AuthMode        string                 `json:"authMode"`
	Config          map[string]interface{} `json:"config"`
	DeployedVersion string                 `json:"deployedVersion"`
	LastDeployedAt  string                 `json:"lastDeployedAt"`
}

// RotateKeyRequest is the body for POST /services/{id}/rotate-key
type RotateKeyRequest struct {
	Environment        string `json:"environment"`
	NotifyOwner        bool   `json:"notifyOwner"`
	GracePeriodMinutes int    `json:"gracePeriodMinutes"`
}

// RotateKeyResponse is returned by POST /services/{id}/rotate-key
type RotateKeyResponse struct {
	OperationID       string `json:"operationId"`
	ServiceID         string `json:"serviceId"`
	Environment       string `json:"environment"`
	KeyID             string `json:"keyId"`
	PreviousKeyID     string `json:"previousKeyId"`
	GracePeriodEndsAt string `json:"gracePeriodEndsAt"`
	ExpiresAt         string `json:"expiresAt"`
	NotificationSent  bool   `json:"notificationSent"`
	Status            string `json:"status"`
}

// ErrorResponse is returned for all errors
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

