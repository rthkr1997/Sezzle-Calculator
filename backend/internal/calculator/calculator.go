package calculator

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidExpression = errors.New("invalid expression")
	ErrDivisionByZero    = errors.New("division by zero")
	ErrNegativeSqrt      = errors.New("cannot take square root of a negative number")
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

func Calculate(expression string) (float64, error) {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return 0, fmt.Errorf("%w: expression is required", ErrInvalidExpression)
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
		return 0, fmt.Errorf("%w: at least one value is required", ErrInvalidExpression)
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
			return 0, fmt.Errorf("%w: power requires exactly two values", ErrInvalidExpression)
		}
		return math.Pow(values[0], values[1]), nil
	case "sqrt":
		if len(values) != 1 {
			return 0, fmt.Errorf("%w: sqrt requires exactly one value", ErrInvalidExpression)
		}
		if values[0] < 0 {
			return 0, ErrNegativeSqrt
		}
		return math.Sqrt(values[0]), nil
	case "percentage":
		if len(values) != 1 {
			return 0, fmt.Errorf("%w: percentage requires exactly one value", ErrInvalidExpression)
		}
		return values[0] / 100, nil
	default:
		return 0, fmt.Errorf("%w: unsupported operation %q", ErrInvalidExpression, op)
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
					return nil, fmt.Errorf("%w: malformed number", ErrInvalidExpression)
				}
				i++
			}
			lit := input[start:i]
			if lit == "." {
				return nil, fmt.Errorf("%w: malformed number", ErrInvalidExpression)
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
				return nil, fmt.Errorf("%w: unknown function %q", ErrInvalidExpression, name)
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
				return nil, fmt.Errorf("%w: unexpected percentage", ErrInvalidExpression)
			}
			tokens = append(tokens, token{operator, "%"})
			expectingOperand = false
		case '+', '*', '/', '^':
			if expectingOperand {
				return nil, fmt.Errorf("%w: unexpected operator %q", ErrInvalidExpression, input[i])
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
			return nil, fmt.Errorf("%w: unsupported character %q", ErrInvalidExpression, input[i])
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
				return nil, fmt.Errorf("%w: mismatched parentheses", ErrInvalidExpression)
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
			return nil, fmt.Errorf("%w: mismatched parentheses", ErrInvalidExpression)
		}
		output = append(output, top)
	}
	return output, nil
}

func evalRPN(tokens []token) (float64, error) {
	var stack []float64
	pop := func() (float64, error) {
		if len(stack) == 0 {
			return 0, fmt.Errorf("%w: missing operand", ErrInvalidExpression)
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
				return 0, fmt.Errorf("%w: malformed number", ErrInvalidExpression)
			}
			stack = append(stack, v)
		case function:
			v, err := pop()
			if err != nil {
				return 0, err
			}
			if v < 0 {
				return 0, ErrNegativeSqrt
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
					return 0, ErrDivisionByZero
				}
				stack = append(stack, left/right)
			case "^":
				stack = append(stack, math.Pow(left, right))
			default:
				return 0, fmt.Errorf("%w: unknown operator", ErrInvalidExpression)
			}
		}
	}
	if len(stack) != 1 {
		return 0, fmt.Errorf("%w: unresolved operands", ErrInvalidExpression)
	}
	return stack[0], nil
}
func isRightAssociative(op string) bool { return op == "^" || op == "u-" }
