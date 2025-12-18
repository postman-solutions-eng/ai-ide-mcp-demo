package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	apphttp "github.com/your-org/ai-ide-mcp-demo/internal/http"
	"github.com/your-org/ai-ide-mcp-demo/internal/store"
)

func main() {
	// Initialize store with seed data
	s := store.New()

	// Initialize handlers
	h := apphttp.NewHandler(s)

	// Create router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health endpoint (no request ID required, no trace headers)
	r.Get("/health", h.Health)

	// API routes (require X-Platform-Request-ID, add trace headers)
	r.Group(func(r chi.Router) {
		r.Use(apphttp.RequireRequestID)
		r.Use(apphttp.AddTraceHeaders)
		r.Use(apphttp.JSONContentType)

		// Services
		r.Get("/services", h.ListServices)
		r.Post("/services", h.CreateService)
		r.Get("/services/{id}", h.GetService)
		r.Patch("/services/{id}", h.UpdateService)

		// Environments
		r.Get("/services/{id}/environments/{env}", h.GetEnvironment)

		// Endpoints
		r.Get("/services/{id}/endpoints", h.GetEndpoints)

		// Actions
		r.Post("/services/{id}/rotate-key", h.RotateKey)
	})

	// Start server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

