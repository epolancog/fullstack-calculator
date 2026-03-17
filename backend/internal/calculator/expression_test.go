package calculator

import (
	"math"
	"testing"
)

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestExpressionPrecedence(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want float64
	}{
		{"multiply before add", "5 + 3 * 2", 11},
		{"mixed precedence", "10 - 2 * 3 + 4", 8},
		{"power before add", "2 ^ 3 + 1", 9},
		{"divide and multiply before add", "10 / 2 + 3 * 4", 17},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluateExpression(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !approxEqual(got, tt.want) {
				t.Errorf("evaluateExpression(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

func TestExpressionBasicOps(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want float64
	}{
		{"addition", "1 + 1", 2},
		{"subtraction", "10 - 3", 7},
		{"multiplication", "4 * 5", 20},
		{"division", "20 / 4", 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluateExpression(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !approxEqual(got, tt.want) {
				t.Errorf("evaluateExpression(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

func TestExpressionEdgeCases(t *testing.T) {
	t.Run("single number", func(t *testing.T) {
		got, err := evaluateExpression("42")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 42 {
			t.Errorf("got %v, want 42", got)
		}
	})

	t.Run("decimal numbers", func(t *testing.T) {
		got, err := evaluateExpression("1.5 + 2.5")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 4 {
			t.Errorf("got %v, want 4", got)
		}
	})

	t.Run("division by zero in expression", func(t *testing.T) {
		_, err := evaluateExpression("10 / 0")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty expression", func(t *testing.T) {
		_, err := evaluateExpression("")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("whitespace only", func(t *testing.T) {
		_, err := evaluateExpression("   ")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid characters", func(t *testing.T) {
		_, err := evaluateExpression("5 & 3")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("consecutive operators", func(t *testing.T) {
		_, err := evaluateExpression("5 + + 3")
		if err == nil {
			t.Fatal("expected error for consecutive operators, got nil")
		}
	})
}

func TestExpressionUnaryMinus(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want float64
	}{
		{"negation at start", "-5 + 3", -2},
		{"negation after operator", "5 + -3", 2},
		{"negation with precedence", "5 * -2 + 1", -9},
		{"double negation", "-5 * -2", 10},
		{"negation with division", "10 / -2", -5},
		{"negation after open paren", "(-5) + 3", -2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluateExpression(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !approxEqual(got, tt.want) {
				t.Errorf("evaluateExpression(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

func TestExpressionPercentage(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want float64
	}{
		{"standalone", "50%", 0.5},
		{"with multiply", "200 * 10%", 20},
		{"with add", "50 + 10%", 55},
		{"with subtract", "50 - 20%", 40},
		{"chained add percent", "100 + 10% + 20%", 132},
		{"chained percent", "50%%", 0.005},
		{"after parens", "(50 + 10)%", 0.6},
		{"mid-precedence", "2 * 3 + 50%", 9},
		{"with power google-style", "2 ^ 3%", math.Pow(2, 0.03)},
		{"standalone percent no left operand", "%", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluateExpression(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !approxEqual(got, tt.want) {
				t.Errorf("evaluateExpression(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

func TestExpressionParentheses(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want float64
	}{
		{"basic grouping", "(5 + 3) * 2", 16},
		{"nested parens", "((2 + 3)) * 4", 20},
		{"right grouping", "10 * (2 + 3)", 50},
		{"complex grouping", "(2 + 3) * (4 - 1)", 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluateExpression(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !approxEqual(got, tt.want) {
				t.Errorf("evaluateExpression(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}

	t.Run("mismatched opening", func(t *testing.T) {
		_, err := evaluateExpression("(5 + 3")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("mismatched closing", func(t *testing.T) {
		_, err := evaluateExpression("5 + 3)")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestExpressionComplex(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want float64
	}{
		{"mixed ops", "2 + 3 * 4 - 6 / 2", 11},
		{"power with multiply", "2 ^ 3 * 2 + 1", 17},
		{"grouped multiply", "(2 + 3) * (4 - 1)", 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluateExpression(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !approxEqual(got, tt.want) {
				t.Errorf("evaluateExpression(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

func TestExpressionSqrt(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want float64
	}{
		{"basic sqrt", "sqrt 16", 4},
		{"sqrt with addition", "sqrt 16 + 9", 13},
		{"sqrt with parens", "sqrt (16 + 9)", 5},
		{"sqrt zero", "sqrt 0", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluateExpression(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !approxEqual(got, tt.want) {
				t.Errorf("evaluateExpression(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}

	t.Run("sqrt negative", func(t *testing.T) {
		_, err := evaluateExpression("sqrt -4")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("sqrt no space rejected", func(t *testing.T) {
		_, err := evaluateExpression("sqrt16")
		if err == nil {
			t.Fatal("expected error for sqrt16 (no space), got nil")
		}
	})
}

func TestExpressionPercentNoLeftOperand(t *testing.T) {
	t.Run("percent 50", func(t *testing.T) {
		got, err := evaluateExpression("% 50")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 0 {
			t.Errorf("evaluateExpression(%q) = %v, want 0", "% 50", got)
		}
	})
}
