package calculator

import (
	"testing"
)

func TestCalculate(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		a        float64
		op       string
		b        float64
		want     float64
		wantErr  bool
	}{
		{"addition", 5, "+", 3, 8, false},
		{"subtraction", 10, "-", 4, 6, false},
		{"multiplication", 6, "*", 8, 48, false},
		{"division", 15, "/", 3, 5, false},
		{"power", 2, "^", 10, 1024, false},
		{"sqrt", 144, "sqrt", 0, 12, false},
		{"percentage", 50, "%", 0, 0.5, false},
		{"division by zero", 10, "/", 0, 0, true},
		{"sqrt negative", -4, "sqrt", 0, 0, true},
		{"unknown operator", 5, "mod", 3, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calc.Calculate(tt.a, tt.op, tt.b)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Calculate(%v, %q, %v) = %v, want %v", tt.a, tt.op, tt.b, got, tt.want)
			}
		})
	}
}

func TestSupportedOperations(t *testing.T) {
	calc := NewCalculator()
	ops := calc.SupportedOperations()

	expected := []string{"%", "*", "+", "-", "/", "^", "sqrt"}
	if len(ops) != len(expected) {
		t.Fatalf("SupportedOperations() returned %d ops, want %d", len(ops), len(expected))
	}
	for i, op := range ops {
		if op != expected[i] {
			t.Errorf("SupportedOperations()[%d] = %q, want %q", i, op, expected[i])
		}
	}
}
