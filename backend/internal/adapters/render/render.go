package render

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func Response(w http.ResponseWriter, statusCode int, contentType string, body []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	if _, err := w.Write(body); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}

func JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode json", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Success(w http.ResponseWriter, statusCode int, message string, data any) {
	JSON(w, statusCode, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, APIResponse{
		Status:  "error",
		Message: message,
	})
}
