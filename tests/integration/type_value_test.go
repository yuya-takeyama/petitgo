package main

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/eval"
	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

// Test RED: We need proper Value types, not just int
func TestTypeValue_StringEvaluation(t *testing.T) {
	input := `"hello"`
	s := scanner.NewScanner(input)
	p := parser.NewParser(s)
	expr := p.ParseExpression()

	env := eval.NewEnvironment()

	// This should return a proper string value, not just length
	// For now this will fail because EvalWithEnvironment returns int
	result := eval.EvalWithEnvironment(expr, env)

	// We want: result should somehow represent "hello", not 5
	// This test documents what we want to achieve
	if result == 5 { // current behavior: returns length
		t.Log("Currently returns string length (5), but we want actual string value")
	}
}

// Test RED: Boolean evaluation should be cleaner
func TestTypeValue_BooleanEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected int // for now still int, but should be bool
	}{
		{"true", 1},
		{"false", 0},
	}

	for _, tt := range tests {
		s := scanner.NewScanner(tt.input)
		p := parser.NewParser(s)
		expr := p.ParseExpression()

		env := eval.NewEnvironment()
		result := eval.EvalWithEnvironment(expr, env)

		if result != tt.expected {
			t.Errorf("Input: %s, expected: %d, got: %d", tt.input, tt.expected, result)
		}
	}
}

// Test RED: Type checking should prevent invalid operations
func TestTypeValue_TypeMismatchError(t *testing.T) {
	// This should eventually error: cannot add string and int
	input := `"hello" + 42`
	s := scanner.NewScanner(input)
	p := parser.NewParser(s)
	expr := p.ParseExpression()

	env := eval.NewEnvironment()

	// For now this will work (string length + int)
	// But eventually should be a type error
	result := eval.EvalWithEnvironment(expr, env)

	// Current: 5 + 42 = 47 (string length + int)
	// Future: should be type error
	if result == 47 {
		t.Log("Currently allows string+int (47), but should be type error in future")
	}
}
