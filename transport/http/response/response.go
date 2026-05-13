package response

import (
	"encoding/json"
	"net/http"
)

// Base adalah struktur balasan standar (Standard Response)
type Base struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// WithJSON mengirimkan balasan sukses berformat JSON
func WithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response := Base{
		Data:    payload,
		Message: "Success",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// WithError mengirimkan balasan error berformat JSON
func WithError(w http.ResponseWriter, err error) {
	response := Base{
		Error: err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest) // Default ke 400 Bad Request
	json.NewEncoder(w).Encode(response)
}