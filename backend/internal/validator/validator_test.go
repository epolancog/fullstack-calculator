package validator

import (
	"errors"
	"math"
	"testing"
)

func newTestValidator() *Validator {
	return NewValidator([]string{"+", "-", "*", "/", "^", "sqrt", "%"})
}

func TestValidateCalculateRequest(t *testing.T) {
	v := newTestValidator()

	t.Run("valid request", func(t *testing.T) {
		if err := v.ValidateCalculateRequest("+", 5, 3); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("empty operator", func(t *testing.T) {
		err := v.ValidateCalculateRequest("", 5, 3)
		assertValidationError(t, err, "INVALID_OPERATOR")
	})

	t.Run("unsupported operator", func(t *testing.T) {
		err := v.ValidateCalculateRequest("mod", 5, 3)
		assertValidationError(t, err, "INVALID_OPERATOR")
	})

	t.Run("NaN operand_a", func(t *testing.T) {
		err := v.ValidateCalculateRequest("+", math.NaN(), 3)
		assertValidationError(t, err, "INVALID_OPERAND")
	})

	t.Run("Inf operand_b", func(t *testing.T) {
		err := v.ValidateCalculateRequest("+", 5, math.Inf(1))
		assertValidationError(t, err, "INVALID_OPERAND")
	})

	t.Run("negative Inf operand_a", func(t *testing.T) {
		err := v.ValidateCalculateRequest("+", math.Inf(-1), 3)
		assertValidationError(t, err, "INVALID_OPERAND")
	})
}

func TestValidateExpressionRequest(t *testing.T) {
	v := newTestValidator()

	t.Run("valid expression", func(t *testing.T) {
		if err := v.ValidateExpressionRequest("5 + 3 * 2"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("empty expression", func(t *testing.T) {
		err := v.ValidateExpressionRequest("")
		assertValidationError(t, err, "INVALID_EXPRESSION")
	})

	t.Run("whitespace only expression", func(t *testing.T) {
		err := v.ValidateExpressionRequest("   ")
		assertValidationError(t, err, "INVALID_EXPRESSION")
	})
}

func assertValidationError(t *testing.T, err error, expectedCode string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ve ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T: %v", err, err)
	}
	if ve.Code != expectedCode {
		t.Errorf("expected code %q, got %q", expectedCode, ve.Code)
	}
}
