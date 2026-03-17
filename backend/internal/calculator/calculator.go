package calculator

import (
	"sort"
)

// Calculator defines the interface for the calculation engine.
type Calculator interface {
	Calculate(operandA float64, operator string, operandB float64) (float64, error)
	EvaluateExpression(expression string) (float64, error)
	SupportedOperations() []string
}

// Calc is the concrete implementation of Calculator backed by an operator registry.
type Calc struct {
	operators map[string]Operator
}

// NewCalculator creates a Calc with all supported operators registered.
func NewCalculator() *Calc {
	return &Calc{
		operators: map[string]Operator{
			"+":    Add{},
			"-":    Subtract{},
			"*":    Multiply{},
			"/":    Divide{},
			"^":    Power{},
			"sqrt": SquareRoot{},
			"%":    Percentage{},
		},
	}
}

// Calculate performs a single binary (or unary) operation.
func (c *Calc) Calculate(operandA float64, operator string, operandB float64) (float64, error) {
	op, ok := c.operators[operator]
	if !ok {
		return 0, InvalidOperatorError{Operator: operator}
	}
	return op.Execute(operandA, operandB)
}

// SupportedOperations returns a sorted list of registered operator keys.
func (c *Calc) SupportedOperations() []string {
	ops := make([]string, 0, len(c.operators))
	for k := range c.operators {
		ops = append(ops, k)
	}
	sort.Strings(ops)
	return ops
}

// EvaluateExpression parses and evaluates a mathematical expression string.
// Implementation is in expression.go.
func (c *Calc) EvaluateExpression(expression string) (float64, error) {
	return evaluateExpression(expression)
}
