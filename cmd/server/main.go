package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	handler "github.com/postman-solutions-eng/ai-ide-mcp-demo/internal/http"
	"github.com/postman-solutions-eng/ai-ide-mcp-demo/internal/store"
)

func main() {
	s := store.New()
	h := handler.NewHandler(s)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health endpoint - no auth required
	r.Get("/health", h.Health)

	// All other endpoints require X-Platform-Request-ID
	r.Group(func(r chi.Router) {
		r.Use(handler.JSONContentType)
		r.Use(handler.RequireRequestID)

		// Services
		r.Get("/services", h.ListServices)
		r.Post("/services", h.CreateService)
		r.Get("/services/{id}", h.GetService)
		r.Patch("/services/{id}", h.UpdateService)
		r.Delete("/services/{id}", h.DeleteService)

		// Environments
		r.Get("/services/{id}/environments/{env}", h.GetEnvironment)

		// Actions
		r.Post("/services/{id}/rotate-key", h.RotateKey)
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

