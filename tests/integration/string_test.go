package main

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/eval"
	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

// Test RED: String parsing should create StringNode, not fail
func TestString_BasicStringLiteral(t *testing.T) {
	input := `"hello world"`
	s := scanner.NewScanner(input)
	p := parser.NewParser(s)

	expr := p.ParseExpression()

	// This should be a StringNode
	stringNode, ok := expr.(*ast.StringNode)
	if !ok {
		t.Fatalf("Expected StringNode, got %T", expr)
	}

	if stringNode.Value != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", stringNode.Value)
	}
}

// Test RED: String with escape sequences
func TestString_EscapeSequences(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello\nworld"`, "hello\nworld"},
		{`"say \"hi\""`, `say "hi"`},
		{`"tab\there"`, "tab\there"},
	}

	for _, tt := range tests {
		s := scanner.NewScanner(tt.input)
		p := parser.NewParser(s)
		expr := p.ParseExpression()

		stringNode, ok := expr.(*ast.StringNode)
		if !ok {
			t.Fatalf("Input: %s, expected StringNode, got %T", tt.input, expr)
		}

		if stringNode.Value != tt.expected {
			t.Errorf("Input: %s, expected: %s, got: %s", tt.input, tt.expected, stringNode.Value)
		}
	}
}

// Test RED: String evaluation should work
func TestString_Evaluation(t *testing.T) {
	input := `"test string"`
	s := scanner.NewScanner(input)
	p := parser.NewParser(s)
	expr := p.ParseExpression()

	env := eval.NewEnvironment()

	// For now we'll have strings return some kind of representation
	// This will fail until we implement proper string evaluation
	result := eval.EvalWithEnvironment(expr, env)

	// We need to define what string evaluation returns
	// For now, let's say it should not return 0 (default error value)
	if result == 0 {
		t.Error("String evaluation should not return 0")
	}
}
