package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// ErrorResponse represents a structured error response for the client.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// JSON writes a JSON response with the given status code and data.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}

// Error writes a structured error response in JSON format.
func Error(w http.ResponseWriter, status int, message string) {
	// Log internal details if needed, but don't send to client
	slog.Warn("handler error", "status", status, "message", message)
	JSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Code:    status,
	})
}
