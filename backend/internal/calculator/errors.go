package calculator

// DivisionByZeroError is returned when dividing by zero.
type DivisionByZeroError struct{}

func (DivisionByZeroError) Error() string {
	return "division by zero is not allowed"
}

// SqrtNegativeError is returned when taking the square root of a negative number.
type SqrtNegativeError struct{}

func (SqrtNegativeError) Error() string {
	return "square root of negative number is not allowed"
}

// InvalidExpressionError is returned for malformed or invalid expressions.
type InvalidExpressionError struct {
	Detail string
}

func (e InvalidExpressionError) Error() string {
	return e.Detail
}

// InvalidOperatorError is returned when an unknown operator is used.
type InvalidOperatorError struct {
	Operator string
}

func (e InvalidOperatorError) Error() string {
	return "unknown operator: " + e.Operator
}
