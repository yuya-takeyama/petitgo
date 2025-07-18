package main

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/eval"
	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

func TestMultipleStatements(t *testing.T) {
	statements := []string{
		"var x int = 10",
		"var y int = 20",
		"z := x + y",
	}

	env := eval.NewEnvironment()

	for _, stmt := range statements {
		sc := scanner.NewScanner(stmt)
		parser := parser.NewParser(sc)
		statement := parser.ParseStatement()
		eval.EvalStatement(statement, env)
	}

	// x = 10
	x, exists := env.GetInt("x")
	if !exists {
		t.Fatalf("expected variable 'x' to exist")
	}
	if x != 10 {
		t.Errorf("expected x = 10, got %d", x)
	}

	// y = 20
	y, exists := env.GetInt("y")
	if !exists {
		t.Fatalf("expected variable 'y' to exist")
	}
	if y != 20 {
		t.Errorf("expected y = 20, got %d", y)
	}

	// z = 30
	z, exists := env.GetInt("z")
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

	env := eval.NewEnvironment()

	for _, stmt := range statements {
		sc := scanner.NewScanner(stmt)
		parser := parser.NewParser(sc)
		statement := parser.ParseStatement()
		eval.EvalStatement(statement, env)
	}

	// x should be 30 ((10 + 5) * 2)
	x, exists := env.GetInt("x")
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

	env := eval.NewEnvironment()

	for _, stmt := range statements {
		sc := scanner.NewScanner(stmt)
		parser := parser.NewParser(sc)
		statement := parser.ParseStatement()
		eval.EvalStatement(statement, env)
	}

	// result should be 14 (2 + 3 * 4)
	result, exists := env.GetInt("result")
	if !exists {
		t.Fatalf("expected variable 'result' to exist")
	}
	if result != 14 {
		t.Errorf("expected result = 14, got %d", result)
	}
}
