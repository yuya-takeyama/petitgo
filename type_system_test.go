package main

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/eval"
	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

// Test RED: String literals should be recognized and parsed
func TestTypeSystem_StringLiterals(t *testing.T) {
	input := `"hello world"`
	s := scanner.NewScanner(input)
	p := parser.NewParser(s)

	expr := p.ParseExpression()

	// This should be a StringNode, not NumberNode
	// Currently this will fail because we don't have STRING token support
	env := eval.NewEnvironment()
	result := eval.EvalWithEnvironment(expr, env)

	// For now, let's expect this to fail gracefully
	// Later we'll implement proper string handling
	_ = result
}

// Test RED: Boolean literals should work
func TestTypeSystem_BooleanLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		s := scanner.NewScanner(tt.input)
		p := parser.NewParser(s)

		expr := p.ParseExpression()

		// This will definitely fail because we don't have boolean support yet
		env := eval.NewEnvironment()
		result := eval.EvalWithEnvironment(expr, env)

		// Convert int result back to bool (1 = true, 0 = false)
		boolResult := result != 0
		if boolResult != tt.expected {
			t.Errorf("Input: %s, expected: %v, got: %v", tt.input, tt.expected, boolResult)
		}
	}
}

// Test RED: String variables should work
func TestTypeSystem_StringVariables(t *testing.T) {
	env := eval.NewEnvironment()

	// Define string variable: var message string = "hello"
	funcDefInput := `var message string = "hello"`
	s1 := scanner.NewScanner(funcDefInput)
	p1 := parser.NewParser(s1)
	stmt := p1.ParseStatement()
	eval.EvalStatement(stmt, env)

	// Use the variable: message
	callInput := "message"
	s2 := scanner.NewScanner(callInput)
	p2 := parser.NewParser(s2)
	expr := p2.ParseExpression()

	result := eval.EvalWithEnvironment(expr, env)

	// This will fail because we don't support strings yet
	// For now, just check it doesn't crash
	_ = result
}

// Test RED: Type checking should prevent mixing types
func TestTypeSystem_TypeChecking(t *testing.T) {
	env := eval.NewEnvironment()

	// This should eventually fail type checking: "hello" + 42
	input := `"hello" + 42`
	s := scanner.NewScanner(input)
	p := parser.NewParser(s)
	expr := p.ParseExpression()

	// For now this will parse but should eventually fail at type checking
	result := eval.EvalWithEnvironment(expr, env)
	_ = result

	// In the future, this should trigger a type error
}
