package handler

import (
	"net/http"

	"my-go-server/pkg/response"
)

// Home serves the home page
func Home(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"message": "Welcome to Go Server"})
}
