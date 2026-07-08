package handler

import (
	"net/http"
	"strconv"

	"github.com/CodeIter/gowebserver/pkg/random"
	"github.com/CodeIter/gowebserver/pkg/response"
)

// RandomPassword generates a random password of a specified length
func RandomPassword(w http.ResponseWriter, r *http.Request) {
	lengthStr := r.PathValue("length")
	chars := r.URL.Query().Get("chars")

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid length value")
		return
	}

	password, err := random.GenerateRandomPassword(length, chars)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "failed to generate random password: "+err.Error())
		return
	}

	response.JSONUnscaped(w, http.StatusOK, map[string]any{"password": password})
}
