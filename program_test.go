package main

import "testing"

func TestMultipleStatements(t *testing.T) {
	statements := []string{
		"var x int = 10",
		"var y int = 20",
		"z := x + y",
	}

	env := NewEnvironment()

	for _, stmt := range statements {
		lexer := NewLexer(stmt)
		parser := NewParser(lexer)
		statement := parser.ParseStatement()
		EvalStatement(statement, env)
	}

	// x = 10
	x, exists := env.Get("x")
	if !exists {
		t.Fatalf("expected variable 'x' to exist")
	}
	if x != 10 {
		t.Errorf("expected x = 10, got %d", x)
	}

	// y = 20
	y, exists := env.Get("y")
	if !exists {
		t.Fatalf("expected variable 'y' to exist")
	}
	if y != 20 {
		t.Errorf("expected y = 20, got %d", y)
	}

	// z = 30
	z, exists := env.Get("z")
	if !exists {
		t.Fatalf("expected variable 'z' to exist")
	}
	if z != 30 {
		t.Errorf("expected z = 30, got %d", z)
	}
}

func TestVariableReassignment(t *testing.T) {
	statements := []string{
		"var x int = 10",
		"x = x + 5",
		"x = x * 2",
	}

	env := NewEnvironment()

	for _, stmt := range statements {
		lexer := NewLexer(stmt)
		parser := NewParser(lexer)
		statement := parser.ParseStatement()
		EvalStatement(statement, env)
	}

	// x should be 30 ((10 + 5) * 2)
	x, exists := env.Get("x")
	if !exists {
		t.Fatalf("expected variable 'x' to exist")
	}
	if x != 30 {
		t.Errorf("expected x = 30, got %d", x)
	}
}

func TestComplexExpressionWithVariables(t *testing.T) {
	statements := []string{
		"var a int = 2",
		"var b int = 3",
		"var c int = 4",
		"result := a + b * c",
	}

	env := NewEnvironment()

	for _, stmt := range statements {
		lexer := NewLexer(stmt)
		parser := NewParser(lexer)
		statement := parser.ParseStatement()
		EvalStatement(statement, env)
	}

	// result should be 14 (2 + 3 * 4)
	result, exists := env.Get("result")
	if !exists {
		t.Fatalf("expected variable 'result' to exist")
	}
	if result != 14 {
		t.Errorf("expected result = 14, got %d", result)
	}
}
