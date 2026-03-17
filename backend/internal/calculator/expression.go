package calculator

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// token represents a parsed token from an expression.
type token struct {
	kind  tokenKind
	value string
}

type tokenKind int

const (
	tokenNumber tokenKind = iota
	tokenOperator
	tokenLeftParen
	tokenRightParen
	tokenFunction // sqrt
	tokenPercent  // % (unary postfix)
)

// operatorInfo holds precedence and associativity for shunting-yard.
type operatorInfo struct {
	precedence     int
	rightAssoc     bool
	unaryPrefix    bool
}

var operators = map[string]operatorInfo{
	"+":    {precedence: 1, rightAssoc: false},
	"-":    {precedence: 1, rightAssoc: false},
	"*":    {precedence: 2, rightAssoc: false},
	"/":    {precedence: 2, rightAssoc: false},
	"^":    {precedence: 3, rightAssoc: true},
	"sqrt": {precedence: 4, unaryPrefix: true},
}

// tokenize splits an expression string into tokens, handling unary minus.
func tokenize(expr string) ([]token, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, errors.New("empty expression")
	}

	var tokens []token
	i := 0

	for i < len(expr) {
		ch := rune(expr[i])

		// Skip whitespace
		if unicode.IsSpace(ch) {
			i++
			continue
		}

		// Check for "sqrt"
		if i+4 <= len(expr) && expr[i:i+4] == "sqrt" {
			// Ensure sqrt is not followed by a digit/letter (require space or paren)
			if i+4 < len(expr) && (unicode.IsDigit(rune(expr[i+4])) || unicode.IsLetter(rune(expr[i+4]))) {
				return nil, fmt.Errorf("invalid token at position %d: sqrt requires a space before its operand", i)
			}
			tokens = append(tokens, token{kind: tokenFunction, value: "sqrt"})
			i += 4
			continue
		}

		// Check for '%'
		if ch == '%' {
			tokens = append(tokens, token{kind: tokenPercent, value: "%"})
			i++
			continue
		}

		// Parentheses
		if ch == '(' {
			tokens = append(tokens, token{kind: tokenLeftParen, value: "("})
			i++
			continue
		}
		if ch == ')' {
			tokens = append(tokens, token{kind: tokenRightParen, value: ")"})
			i++
			continue
		}

		// Number (including unary minus)
		if unicode.IsDigit(ch) || ch == '.' || (ch == '-' && isUnaryMinus(tokens)) {
			start := i
			if ch == '-' {
				i++ // consume the minus sign
			}
			hasDecimal := false
			for i < len(expr) && (unicode.IsDigit(rune(expr[i])) || rune(expr[i]) == '.') {
				if rune(expr[i]) == '.' {
					if hasDecimal {
						return nil, fmt.Errorf("invalid number at position %d: multiple decimal points", start)
					}
					hasDecimal = true
				}
				i++
			}
			numStr := expr[start:i]
			if numStr == "-" || numStr == "." || numStr == "-." {
				return nil, fmt.Errorf("invalid number at position %d: %q", start, numStr)
			}
			tokens = append(tokens, token{kind: tokenNumber, value: numStr})
			continue
		}

		// Operators: +, -, *, /, ^
		if ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '^' {
			tokens = append(tokens, token{kind: tokenOperator, value: string(ch)})
			i++
			continue
		}

		return nil, fmt.Errorf("invalid character at position %d: %q", i, string(ch))
	}

	return tokens, nil
}

// isUnaryMinus returns true if a '-' at the current position should be treated as unary negation.
func isUnaryMinus(preceding []token) bool {
	if len(preceding) == 0 {
		return true
	}
	last := preceding[len(preceding)-1]
	return last.kind == tokenOperator || last.kind == tokenLeftParen || last.kind == tokenFunction
}

// rpnToken is a token in the RPN (postfix) output queue.
type rpnToken struct {
	kind  tokenKind
	value string
}

// shuntingYard converts infix tokens to RPN (postfix) using Dijkstra's shunting-yard algorithm.
func shuntingYard(tokens []token) ([]rpnToken, error) {
	var output []rpnToken
	var opStack []token

	for _, tok := range tokens {
		switch tok.kind {
		case tokenNumber:
			output = append(output, rpnToken{kind: tokenNumber, value: tok.value})

		case tokenFunction:
			opStack = append(opStack, tok)

		case tokenPercent:
			// Postfix unary — goes directly to output
			output = append(output, rpnToken{kind: tokenPercent, value: "%"})

		case tokenOperator:
			info := operators[tok.value]
			for len(opStack) > 0 {
				top := opStack[len(opStack)-1]
				if top.kind == tokenLeftParen {
					break
				}
				topInfo, isOp := operators[top.value]
				if !isOp {
					break
				}
				if topInfo.precedence > info.precedence ||
					(topInfo.precedence == info.precedence && !info.rightAssoc) {
					output = append(output, rpnToken{kind: top.kind, value: top.value})
					opStack = opStack[:len(opStack)-1]
				} else {
					break
				}
			}
			opStack = append(opStack, tok)

		case tokenLeftParen:
			opStack = append(opStack, tok)

		case tokenRightParen:
			found := false
			for len(opStack) > 0 {
				top := opStack[len(opStack)-1]
				opStack = opStack[:len(opStack)-1]
				if top.kind == tokenLeftParen {
					found = true
					break
				}
				output = append(output, rpnToken{kind: top.kind, value: top.value})
			}
			if !found {
				return nil, errors.New("mismatched parentheses: extra closing parenthesis")
			}
			// If the top of the stack is a function, pop it to output
			if len(opStack) > 0 && opStack[len(opStack)-1].kind == tokenFunction {
				top := opStack[len(opStack)-1]
				opStack = opStack[:len(opStack)-1]
				output = append(output, rpnToken{kind: top.kind, value: top.value})
			}
		}
	}

	// Pop remaining operators
	for len(opStack) > 0 {
		top := opStack[len(opStack)-1]
		opStack = opStack[:len(opStack)-1]
		if top.kind == tokenLeftParen {
			return nil, errors.New("mismatched parentheses: extra opening parenthesis")
		}
		output = append(output, rpnToken{kind: top.kind, value: top.value})
	}

	return output, nil
}

// evaluateRPN evaluates an RPN token queue, handling context-dependent percentage.
func evaluateRPN(rpn []rpnToken) (float64, error) {
	var stack []float64
	// Track the last binary operator pushed for context-dependent %
	var lastBinaryOp string
	// Track if % was applied with no left operand (implicit 0)
	implicitZeroPercent := false

	for idx, tok := range rpn {
		switch tok.kind {
		case tokenNumber:
			val, err := strconv.ParseFloat(tok.value, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %q", tok.value)
			}
			stack = append(stack, val)

		case tokenPercent:
			if len(stack) == 0 {
				// No left operand: implicit 0 → result is 0
				stack = append(stack, 0)
				implicitZeroPercent = true
				continue
			}
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			// Determine the pending binary operator to decide % behavior.
			// Look ahead in the RPN to find the next binary operator.
			pendingOp := findPendingBinaryOp(rpn, idx)
			if pendingOp == "" {
				pendingOp = lastBinaryOp
			}

			switch pendingOp {
			case "+", "-":
				// Percent-of-left-operand: we need the left operand of the pending +/-
				// which is the current top of the stack.
				if len(stack) > 0 {
					base := stack[len(stack)-1]
					stack = append(stack, base*(a/100))
				} else {
					// No base to compute percent of — just divide by 100
					stack = append(stack, a/100)
				}
			default:
				// Simple divide by 100 (for *, /, ^, standalone)
				stack = append(stack, a/100)
			}

		case tokenFunction:
			if tok.value == "sqrt" {
				if len(stack) == 0 {
					return 0, errors.New("sqrt requires an operand")
				}
				a := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if a < 0 {
					return 0, errors.New("square root of negative number is not allowed")
				}
				stack = append(stack, math.Sqrt(a))
			}

		case tokenOperator:
			if len(stack) < 2 {
				return 0, fmt.Errorf("not enough operands for operator %q", tok.value)
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch tok.value {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, errors.New("division by zero is not allowed")
				}
				result = a / b
			case "^":
				result = math.Pow(a, b)
			default:
				return 0, fmt.Errorf("unknown operator: %q", tok.value)
			}
			stack = append(stack, result)
			lastBinaryOp = tok.value
		}
	}

	if len(stack) == 0 {
		return 0, errors.New("empty expression")
	}
	if len(stack) > 1 {
		// If % was applied with no left operand (implicit 0), return 0
		if implicitZeroPercent {
			return 0, nil
		}
		return 0, errors.New("malformed expression: too many values")
	}
	return stack[0], nil
}

// findPendingBinaryOp looks ahead in the RPN from the current index
// to find the next binary operator that will consume the % result.
func findPendingBinaryOp(rpn []rpnToken, currentIdx int) string {
	for i := currentIdx + 1; i < len(rpn); i++ {
		if rpn[i].kind == tokenOperator {
			return rpn[i].value
		}
	}
	return ""
}

// evaluateExpression is the main entry point: tokenize → shunting-yard → evaluate RPN.
func evaluateExpression(expr string) (float64, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return 0, err
	}

	rpn, err := shuntingYard(tokens)
	if err != nil {
		return 0, err
	}

	return evaluateRPN(rpn)
}
