package handler

import (
	"net/http"

	"my-go-server/pkg/response"
)

// Ready checks if the server is ready to accept requests
func Ready(w http.ResponseWriter, r *http.Request) {
	// Check DB/Dependencies here
	response.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
