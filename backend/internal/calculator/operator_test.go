package calculator

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	op := Add{}
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed signs", -2, 3, 1},
		{"zeros", 0, 0, 0},
		{"decimals", 1.5, 2.5, 4},
		{"large numbers", 999999, 1, 1000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := op.Execute(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Add.Execute(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	op := Subtract{}
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive numbers", 5, 3, 2},
		{"negative numbers", -2, -3, 1},
		{"result negative", 3, 5, -2},
		{"zeros", 0, 0, 0},
		{"decimals", 5.5, 2.5, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := op.Execute(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Subtract.Execute(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	op := Multiply{}
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive numbers", 4, 5, 20},
		{"negative numbers", -2, -3, 6},
		{"mixed signs", -2, 3, -6},
		{"by zero", 5, 0, 0},
		{"decimals", 1.5, 2, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := op.Execute(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Multiply.Execute(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	op := Divide{}
	t.Run("normal division", func(t *testing.T) {
		got, err := op.Execute(10, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 5 {
			t.Errorf("Divide.Execute(10, 2) = %v, want 5", got)
		}
	})
	t.Run("decimal result", func(t *testing.T) {
		got, err := op.Execute(7, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 3.5 {
			t.Errorf("Divide.Execute(7, 2) = %v, want 3.5", got)
		}
	})
	t.Run("division by zero", func(t *testing.T) {
		_, err := op.Execute(10, 0)
		if err == nil {
			t.Fatal("expected error for division by zero, got nil")
		}
		if err.Error() != "division by zero is not allowed" {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}

func TestPower(t *testing.T) {
	op := Power{}
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive exponent", 2, 3, 8},
		{"zero exponent", 5, 0, 1},
		{"negative exponent", 2, -1, 0.5},
		{"base zero", 0, 5, 0},
		{"one to any power", 1, 100, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := op.Execute(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Power.Execute(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSquareRoot(t *testing.T) {
	op := SquareRoot{}
	t.Run("positive number", func(t *testing.T) {
		got, err := op.Execute(16, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 4 {
			t.Errorf("SquareRoot.Execute(16, 0) = %v, want 4", got)
		}
	})
	t.Run("zero", func(t *testing.T) {
		got, err := op.Execute(0, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 0 {
			t.Errorf("SquareRoot.Execute(0, 0) = %v, want 0", got)
		}
	})
	t.Run("negative number", func(t *testing.T) {
		_, err := op.Execute(-4, 0)
		if err == nil {
			t.Fatal("expected error for sqrt of negative, got nil")
		}
		if err.Error() != "square root of negative number is not allowed" {
			t.Errorf("unexpected error message: %v", err)
		}
	})
	t.Run("perfect square", func(t *testing.T) {
		got, err := op.Execute(144, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 12 {
			t.Errorf("SquareRoot.Execute(144, 0) = %v, want 12", got)
		}
	})
	t.Run("non-perfect square", func(t *testing.T) {
		got, err := op.Execute(2, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if math.Abs(got-math.Sqrt(2)) > 1e-9 {
			t.Errorf("SquareRoot.Execute(2, 0) = %v, want %v", got, math.Sqrt(2))
		}
	})
}

func TestPercentage(t *testing.T) {
	op := Percentage{}
	tests := []struct {
		name string
		a    float64
		want float64
	}{
		{"fifty percent", 50, 0.5},
		{"hundred percent", 100, 1.0},
		{"zero percent", 0, 0},
		{"decimal input", 12.5, 0.125},
		{"negative input", -50, -0.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := op.Execute(tt.a, 0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Percentage.Execute(%v, 0) = %v, want %v", tt.a, got, tt.want)
			}
		})
	}
}
