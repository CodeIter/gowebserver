package handler

import (
	"net/http"
	"strconv"

	"github.com/CodeIter/gowebserver/pkg/random"
	"github.com/CodeIter/gowebserver/pkg/response"
)

// RandomNumber generates a random number within a specified range
func RandomNumber(w http.ResponseWriter, r *http.Request) {
	minStr := r.PathValue("min")
	maxStr := r.PathValue("max")

	min, err := strconv.Atoi(minStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid min value")
		return
	}

	max, err := strconv.Atoi(maxStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid max value")
		return
	}

	num, err := random.GenerateRandomNumber(min, max)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "failed to generate random number: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{"number": num})
}
