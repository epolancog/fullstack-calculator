package calculator

import (
	"errors"
	"math"
)

// Operator defines the interface for arithmetic operations.
type Operator interface {
	Execute(a, b float64) (float64, error)
}

// Add performs addition.
type Add struct{}

func (Add) Execute(a, b float64) (float64, error) {
	return a + b, nil
}

// Subtract performs subtraction.
type Subtract struct{}

func (Subtract) Execute(a, b float64) (float64, error) {
	return a - b, nil
}

// Multiply performs multiplication.
type Multiply struct{}

func (Multiply) Execute(a, b float64) (float64, error) {
	return a * b, nil
}

// Divide performs division, returning an error on division by zero.
type Divide struct{}

func (Divide) Execute(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero is not allowed")
	}
	return a / b, nil
}

// Power performs exponentiation (a^b).
type Power struct{}

func (Power) Execute(a, b float64) (float64, error) {
	return math.Pow(a, b), nil
}

// SquareRoot computes the square root of a. Operand b is ignored.
type SquareRoot struct{}

func (SquareRoot) Execute(a, _ float64) (float64, error) {
	if a < 0 {
		return 0, errors.New("square root of negative number is not allowed")
	}
	return math.Sqrt(a), nil
}

// Percentage divides a by 100. Operand b is ignored.
// Context-dependent behavior (percent-of-left-operand with +/-) is handled
// in the expression evaluator, not here.
type Percentage struct{}

func (Percentage) Execute(a, _ float64) (float64, error) {
	return a / 100, nil
}
