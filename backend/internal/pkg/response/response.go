package response

import (
	"encoding/json"
	"net/http"
)

// APIResponse represents the standard structure for all API responses
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// JSON sends a standardized JSON response
func JSON(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}, errors interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := APIResponse{
		Success: success,
		Message: message,
		Data:    data,
		Errors:  errors,
	}

	json.NewEncoder(w).Encode(resp)
}

// Success is a helper for successful responses
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	JSON(w, statusCode, true, message, data, nil)
}

// Error is a helper for error responses
func Error(w http.ResponseWriter, statusCode int, message string, errors interface{}) {
	JSON(w, statusCode, false, message, nil, errors)
}
