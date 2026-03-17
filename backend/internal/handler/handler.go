package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/epolancog/fullstack-calculator/backend/internal/calculator"
	"github.com/epolancog/fullstack-calculator/backend/internal/validator"
)

// Handler holds dependencies for the HTTP handlers.
type Handler struct {
	calc      calculator.Calculator
	validator *validator.Validator
}

// NewHandler creates a Handler with the given calculator and validator.
func NewHandler(calc calculator.Calculator, v *validator.Validator) *Handler {
	return &Handler{calc: calc, validator: v}
}

// RegisterRoutes registers all API routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/calculate", h.handleCalculate)
	mux.HandleFunc("POST /api/calculate/expression", h.handleExpression)
	mux.HandleFunc("GET /api/operations", h.handleOperations)
}

func (h *Handler) handleCalculate(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "MALFORMED_JSON", "invalid JSON: "+err.Error())
		return
	}

	if err := h.validator.ValidateCalculateRequest(req.Operator, req.OperandA, req.OperandB); err != nil {
		var ve validator.ValidationError
		if errors.As(err, &ve) {
			writeError(w, http.StatusBadRequest, ve.Code, ve.Message)
			return
		}
		writeError(w, http.StatusBadRequest, "INVALID_OPERAND", err.Error())
		return
	}

	result, err := h.calc.Calculate(req.OperandA, req.Operator, req.OperandB)
	if err != nil {
		writeCalcError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CalculateResponse{Result: result})
}

func (h *Handler) handleExpression(w http.ResponseWriter, r *http.Request) {
	var req ExpressionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "MALFORMED_JSON", "invalid JSON: "+err.Error())
		return
	}

	if err := h.validator.ValidateExpressionRequest(req.Expression); err != nil {
		var ve validator.ValidationError
		if errors.As(err, &ve) {
			writeError(w, http.StatusBadRequest, ve.Code, ve.Message)
			return
		}
		writeError(w, http.StatusBadRequest, "INVALID_EXPRESSION", err.Error())
		return
	}

	result, err := h.calc.EvaluateExpression(req.Expression)
	if err != nil {
		writeCalcError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, ExpressionResponse{
		Result:     result,
		Expression: req.Expression,
	})
}

func (h *Handler) handleOperations(w http.ResponseWriter, r *http.Request) {
	ops := h.calc.SupportedOperations()
	writeJSON(w, http.StatusOK, OperationsResponse{Operations: ops})
}

// writeCalcError maps calculator typed errors to HTTP error responses.
func writeCalcError(w http.ResponseWriter, err error) {
	var divErr calculator.DivisionByZeroError
	if errors.As(err, &divErr) {
		writeError(w, http.StatusBadRequest, "DIVISION_BY_ZERO", err.Error())
		return
	}

	var sqrtErr calculator.SqrtNegativeError
	if errors.As(err, &sqrtErr) {
		writeError(w, http.StatusBadRequest, "SQRT_NEGATIVE", err.Error())
		return
	}

	var invOpErr calculator.InvalidOperatorError
	if errors.As(err, &invOpErr) {
		writeError(w, http.StatusBadRequest, "INVALID_OPERATOR", err.Error())
		return
	}

	var invExprErr calculator.InvalidExpressionError
	if errors.As(err, &invExprErr) {
		writeError(w, http.StatusBadRequest, "INVALID_EXPRESSION", err.Error())
		return
	}

	writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "unexpected error")
}
