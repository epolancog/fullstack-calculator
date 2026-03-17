package handler

// CalculateRequest is the request body for single arithmetic operations.
type CalculateRequest struct {
	OperandA float64 `json:"operand_a"`
	Operator string  `json:"operator"`
	OperandB float64 `json:"operand_b"`
}

// ExpressionRequest is the request body for full expression evaluation.
type ExpressionRequest struct {
	Expression string `json:"expression"`
}
