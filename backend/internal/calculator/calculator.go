package calculator

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type ErrorCode string

const (
	ErrorCodeInvalidExpression ErrorCode = "invalid_expression"
	ErrorCodeDivisionByZero    ErrorCode = "division_by_zero_not_allowed"
	ErrorCodeNegativeSqrt      ErrorCode = "square_root_of_negative_number_not_allowed"
)

type CalculatorError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *CalculatorError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return string(e.Code)
}

func (e *CalculatorError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

var (
	ErrInvalidExpression = &CalculatorError{Code: ErrorCodeInvalidExpression, Message: "invalid expression"}
	ErrDivisionByZero    = &CalculatorError{Code: ErrorCodeDivisionByZero, Message: "division by zero is not allowed"}
	ErrNegativeSqrt      = &CalculatorError{Code: ErrorCodeNegativeSqrt, Message: "square root of a negative number is not allowed"}
)

type tokenType int

const (
	number tokenType = iota
	operator
	leftParen
	rightParen
	function
)

type token struct {
	typeOf tokenType
	value  string
}

var precedence = map[string]int{"+": 1, "-": 1, "*": 2, "/": 2, "^": 3, "u-": 4}

func wrapCalculatorError(code ErrorCode, base error, format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	if base == nil {
		base = ErrInvalidExpression
	}
	if code == ErrorCodeInvalidExpression {
		return &CalculatorError{Code: code, Message: fmt.Sprintf("%s: %s", base.Error(), msg), Err: base}
	}
	return &CalculatorError{Code: code, Message: base.Error(), Err: base}
}

func Calculate(expression string) (float64, error) {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return 0, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "expression is required")
	}
	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}
	rpn, err := toRPN(tokens)
	if err != nil {
		return 0, err
	}
	return evalRPN(rpn)
}

func Operation(op string, values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "at least one value is required")
	}
	switch op {
	case "add":
		r := 0.0
		for _, v := range values {
			r += v
		}
		return r, nil
	case "subtract":
		r := values[0]
		for _, v := range values[1:] {
			r -= v
		}
		return r, nil
	case "multiply":
		r := 1.0
		for _, v := range values {
			r *= v
		}
		return r, nil
	case "divide":
		r := values[0]
		for _, v := range values[1:] {
			if v == 0 {
				return 0, ErrDivisionByZero
			}
			r /= v
		}
		return r, nil
	case "power":
		if len(values) != 2 {
			return 0, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "power requires exactly two values")
		}
		return math.Pow(values[0], values[1]), nil
	case "sqrt":
		if len(values) != 1 {
			return 0, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "sqrt requires exactly one value")
		}
		if values[0] < 0 {
			return 0, wrapCalculatorError(ErrorCodeNegativeSqrt, ErrNegativeSqrt, "cannot take square root of a negative number")
		}
		return math.Sqrt(values[0]), nil
	case "percentage":
		if len(values) != 1 {
			return 0, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "percentage requires exactly one value")
		}
		return values[0] / 100, nil
	default:
		return 0, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "unsupported operation %q", op)
	}
}

func tokenize(input string) ([]token, error) {
	var tokens []token
	expectingOperand := true
	for i := 0; i < len(input); {
		ch := rune(input[i])
		if unicode.IsSpace(ch) {
			i++
			continue
		}
		if unicode.IsDigit(ch) || ch == '.' {
			start := i
			dots := 0
			for i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.') {
				if input[i] == '.' {
					dots++
				}
				if dots > 1 {
					return nil, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "malformed number")
				}
				i++
			}
			lit := input[start:i]
			if lit == "." {
				return nil, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "malformed number")
			}
			tokens = append(tokens, token{number, lit})
			expectingOperand = false
			continue
		}
		if unicode.IsLetter(ch) {
			start := i
			for i < len(input) && unicode.IsLetter(rune(input[i])) {
				i++
			}
			name := strings.ToLower(input[start:i])
			if name != "sqrt" {
				return nil, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "unknown function %q", name)
			}
			tokens = append(tokens, token{function, name})
			expectingOperand = true
			continue
		}
		switch input[i] {
		case '(':
			tokens = append(tokens, token{leftParen, "("})
			expectingOperand = true
		case ')':
			tokens = append(tokens, token{rightParen, ")"})
			expectingOperand = false
		case '%':
			if expectingOperand {
				return nil, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "unexpected percentage")
			}
			tokens = append(tokens, token{operator, "%"})
			expectingOperand = false
		case '+', '*', '/', '^':
			if expectingOperand {
				return nil, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "unexpected operator %q", input[i])
			}
			tokens = append(tokens, token{operator, string(input[i])})
			expectingOperand = true
		case '-':
			if expectingOperand {
				tokens = append(tokens, token{operator, "u-"})
			} else {
				tokens = append(tokens, token{operator, "-"})
				expectingOperand = true
			}
		default:
			return nil, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "unsupported character %q", input[i])
		}
		i++
	}
	return tokens, nil
}

func toRPN(tokens []token) ([]token, error) {
	var output, stack []token
	for _, t := range tokens {
		switch t.typeOf {
		case number:
			output = append(output, t)
		case function:
			stack = append(stack, t)
		case operator:
			if t.value == "%" {
				output = append(output, t)
				continue
			}
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top.typeOf == leftParen {
					break
				}
				if top.typeOf == function || precedence[top.value] > precedence[t.value] || (precedence[top.value] == precedence[t.value] && !isRightAssociative(t.value)) {
					output = append(output, top)
					stack = stack[:len(stack)-1]
					continue
				}
				break
			}
			stack = append(stack, t)
		case leftParen:
			stack = append(stack, t)
		case rightParen:
			found := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top.typeOf == leftParen {
					found = true
					break
				}
				output = append(output, top)
			}
			if !found {
				return nil, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "mismatched parentheses")
			}
			if len(stack) > 0 && stack[len(stack)-1].typeOf == function {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
		}
	}
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top.typeOf == leftParen || top.typeOf == rightParen {
			return nil, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "mismatched parentheses")
		}
		output = append(output, top)
	}
	return output, nil
}

func evalRPN(tokens []token) (float64, error) {
	var stack []float64
	pop := func() (float64, error) {
		if len(stack) == 0 {
			return 0, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "missing operand")
		}
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return v, nil
	}
	for _, t := range tokens {
		switch t.typeOf {
		case number:
			v, err := strconv.ParseFloat(t.value, 64)
			if err != nil {
				return 0, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "malformed number")
			}
			stack = append(stack, v)
		case function:
			v, err := pop()
			if err != nil {
				return 0, err
			}
			if v < 0 {
				return 0, wrapCalculatorError(ErrorCodeNegativeSqrt, ErrNegativeSqrt, "cannot take square root of a negative number")
			}
			stack = append(stack, math.Sqrt(v))
		case operator:
			if t.value == "%" {
				v, err := pop()
				if err != nil {
					return 0, err
				}
				stack = append(stack, v/100)
				continue
			}
			if t.value == "u-" {
				v, err := pop()
				if err != nil {
					return 0, err
				}
				stack = append(stack, -v)
				continue
			}
			right, err := pop()
			if err != nil {
				return 0, err
			}
			left, err := pop()
			if err != nil {
				return 0, err
			}
			switch t.value {
			case "+":
				stack = append(stack, left+right)
			case "-":
				stack = append(stack, left-right)
			case "*":
				stack = append(stack, left*right)
			case "/":
				if right == 0 {
					return 0, wrapCalculatorError(ErrorCodeDivisionByZero, ErrDivisionByZero, "division by zero")
				}
				stack = append(stack, left/right)
			case "^":
				stack = append(stack, math.Pow(left, right))
			default:
				return 0, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "unknown operator")
			}
		}
	}
	if len(stack) != 1 {
		return 0, wrapCalculatorError(ErrorCodeInvalidExpression, ErrInvalidExpression, "unresolved operands")
	}
	return stack[0], nil
}
func isRightAssociative(op string) bool { return op == "^" || op == "u-" }
