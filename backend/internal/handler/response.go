package handler

import (
	"encoding/json"
	"net/http"
)

// CalculateResponse is the response body for single arithmetic operations.
type CalculateResponse struct {
	Result float64 `json:"result"`
}

// ExpressionResponse is the response body for expression evaluation.
type ExpressionResponse struct {
	Result     float64 `json:"result"`
	Expression string  `json:"expression"`
}

// OperationsResponse is the response body for listing supported operations.
type OperationsResponse struct {
	Operations []string `json:"operations"`
}

// ErrorResponse is the standard error response body.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains the error code and human-readable message.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError writes a JSON error response with the given status code, error code, and message.
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}
