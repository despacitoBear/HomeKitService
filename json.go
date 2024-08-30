package main

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JsonResponser(w http.ResponseWriter, statusCode int, status string, err error, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := JsonResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
	if err != nil {
		response.Error = err.Error()
	}
	json.NewEncoder(w).Encode(response)
}
