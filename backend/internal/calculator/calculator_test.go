package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculateExpressions(t *testing.T) {
	cases := []struct {
		expr string
		want float64
	}{
		{"1 + 2 * 3", 7},
		{"(1 + 2) * 3", 9},
		{"2 ^ 3 ^ 2", 512},
		{"sqrt(81) + 50%", 9.5},
		{"-5 + 10 / 2", 0},
		{"12 + 4 * (8 - 3) / 2", 22},
	}
	for _, tc := range cases {
		got, err := Calculate(tc.expr)
		if err != nil {
			t.Fatalf("Calculate(%q) returned error: %v", tc.expr, err)
		}
		if math.Abs(got-tc.want) > 1e-9 {
			t.Fatalf("Calculate(%q) = %v, want %v", tc.expr, got, tc.want)
		}
	}
}

func TestCalculateErrors(t *testing.T) {
	cases := []struct {
		expr   string
		target error
	}{{"10 / 0", ErrDivisionByZero}, {"sqrt(-1)", ErrNegativeSqrt}, {"1 +", ErrInvalidExpression}, {"", ErrInvalidExpression}}
	for _, tc := range cases {
		_, err := Calculate(tc.expr)
		if !errors.Is(err, tc.target) {
			t.Fatalf("Calculate(%q) error = %v, want %v", tc.expr, err, tc.target)
		}
	}
}

func TestOperationMultipleValues(t *testing.T) {
	got, err := Operation("add", []float64{1, 2, 3, 4})
	if err != nil || got != 10 {
		t.Fatalf("add got %v err %v", got, err)
	}
	got, err = Operation("divide", []float64{100, 2, 5})
	if err != nil || got != 10 {
		t.Fatalf("divide got %v err %v", got, err)
	}
}

func TestDefinedCalculatorErrors(t *testing.T) {
	_, err := Calculate("")
	if !errors.Is(err, ErrInvalidExpression) {
		t.Fatalf("expected invalid expression error, got %v", err)
	}
	var calcErr *CalculatorError
	if !errors.As(err, &calcErr) {
		t.Fatalf("expected calculator error, got %T", err)
	}
	if calcErr.Code != ErrorCodeInvalidExpression {
		t.Fatalf("expected invalid expression code, got %q", calcErr.Code)
	}
}
