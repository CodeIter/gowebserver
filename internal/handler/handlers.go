package handler

import (
	"my-go-server/pkg/response"
	"net/http"
	"time"
)

// Health checks if the server is healthy
func Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// Ready checks if the server is ready to accept requests
func Ready(w http.ResponseWriter, r *http.Request) {
	// Check DB/Dependencies here
	response.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

// Home serves the home page
func Home(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"message": "Welcome to Go Server"})
}
