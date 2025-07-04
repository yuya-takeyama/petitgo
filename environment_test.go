package main

import "testing"

func TestEnvironmentSetAndGet(t *testing.T) {
	env := NewEnvironment()

	env.Set("x", 42)

	value, exists := env.Get("x")
	if !exists {
		t.Fatalf("expected variable 'x' to exist")
	}

	if value != 42 {
		t.Errorf("expected value 42, got %d", value)
	}
}

func TestEnvironmentGetNonExistent(t *testing.T) {
	env := NewEnvironment()

	_, exists := env.Get("nonexistent")
	if exists {
		t.Fatalf("expected variable 'nonexistent' to not exist")
	}
}

func TestEnvironmentOverride(t *testing.T) {
	env := NewEnvironment()

	env.Set("x", 42)
	env.Set("x", 100)

	value, exists := env.Get("x")
	if !exists {
		t.Fatalf("expected variable 'x' to exist")
	}

	if value != 100 {
		t.Errorf("expected value 100, got %d", value)
	}
}

func TestStatementEvaluation(t *testing.T) {
	tests := []struct {
		statement string
		varName   string
		expected  int
	}{
		{"var x int = 42", "x", 42},
		{"y := 10", "y", 10},
	}

	for _, tt := range tests {
		env := NewEnvironment()
		lexer := NewLexer(tt.statement)
		parser := NewParser(lexer)
		stmt := parser.ParseStatement()

		EvalStatement(stmt, env)

		value, exists := env.Get(tt.varName)
		if !exists {
			t.Fatalf("expected variable '%s' to exist after statement: %s", tt.varName, tt.statement)
		}

		if value != tt.expected {
			t.Errorf("statement: %s, expected: %d, got: %d", tt.statement, tt.expected, value)
		}
	}
}

func TestVariableReference(t *testing.T) {
	env := NewEnvironment()
	env.Set("x", 42)
	env.Set("y", 10)

	// x + y = 52
	input := "x + y"
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	expr := parser.ParseExpression()

	result := EvalWithEnvironment(expr, env)

	if result != 52 {
		t.Errorf("expected 52, got %d", result)
	}
}

func TestAssignmentToExistingVariable(t *testing.T) {
	env := NewEnvironment()

	// var x int = 42
	lexer1 := NewLexer("var x int = 42")
	parser1 := NewParser(lexer1)
	stmt1 := parser1.ParseStatement()
	EvalStatement(stmt1, env)

	// x = 100
	lexer2 := NewLexer("x = 100")
	parser2 := NewParser(lexer2)
	stmt2 := parser2.ParseStatement()
	EvalStatement(stmt2, env)

	value, exists := env.Get("x")
	if !exists {
		t.Fatalf("expected variable 'x' to exist")
	}

	if value != 100 {
		t.Errorf("expected value 100, got %d", value)
	}
}
