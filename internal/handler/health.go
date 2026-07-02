package handler

import (
	"net/http"
	"time"

	"my-go-server/pkg/response"
)

// Health checks if the server is healthy
func Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}
