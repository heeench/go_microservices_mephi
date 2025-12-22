package handlers

import (
	"io"
	"net/http"

	"go-microservice/services"
)

// IntegrationHandler exposes simple MinIO-backed endpoints.
type IntegrationHandler struct {
	Service *services.IntegrationService
}

func (h *IntegrationHandler) RegisterRoutes(r *http.ServeMux) {
	// Leaving as stub if needed for future expansion.
	// Example: r.HandleFunc("/api/upload", h.Upload)
	_ = r
}

// Upload demonstrates sending a payload to MinIO (not wired by default).
func (h *IntegrationHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if h.Service == nil {
		http.Error(w, "integration service not configured", http.StatusServiceUnavailable)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.Upload(r.Context(), "sample-object", data, "application/octet-stream"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
