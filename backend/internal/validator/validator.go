package validator

import (
	"math"
	"strings"
)

// Validator validates incoming API requests.
type Validator struct {
	supportedOps map[string]bool
}

// NewValidator creates a Validator with the given list of supported operator strings.
func NewValidator(ops []string) *Validator {
	m := make(map[string]bool, len(ops))
	for _, op := range ops {
		m[op] = true
	}
	return &Validator{supportedOps: m}
}

// ValidationError represents a validation failure with an error code.
type ValidationError struct {
	Code    string
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// ValidateCalculateRequest validates the inputs for a single arithmetic operation.
func (v *Validator) ValidateCalculateRequest(operator string, operandA, operandB float64) error {
	if operator == "" {
		return ValidationError{Code: "INVALID_OPERATOR", Message: "operator is required"}
	}
	if !v.supportedOps[operator] {
		return ValidationError{Code: "INVALID_OPERATOR", Message: "unsupported operator: " + operator}
	}
	if math.IsNaN(operandA) || math.IsInf(operandA, 0) {
		return ValidationError{Code: "INVALID_OPERAND", Message: "operand_a must be a finite number"}
	}
	if math.IsNaN(operandB) || math.IsInf(operandB, 0) {
		return ValidationError{Code: "INVALID_OPERAND", Message: "operand_b must be a finite number"}
	}
	return nil
}

// ValidateExpressionRequest validates the inputs for expression evaluation.
func (v *Validator) ValidateExpressionRequest(expression string) error {
	if strings.TrimSpace(expression) == "" {
		return ValidationError{Code: "INVALID_EXPRESSION", Message: "expression is required"}
	}
	return nil
}
