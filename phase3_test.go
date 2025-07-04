package main

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/eval"
	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

func TestPhase3_SimpleComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"5 > 3", 1},
		{"3 > 5", 0},
		{"5 == 5", 1},
		{"5 == 3", 0},
		{"5 != 3", 1},
		{"5 != 5", 0},
	}

	for _, tt := range tests {
		sc := scanner.NewScanner(tt.input)
		parser := parser.NewParser(sc)
		ast := parser.ParseExpression()
		result := eval.Eval(ast)

		if result != tt.expected {
			t.Errorf("input: %s, expected: %d, got: %d", tt.input, tt.expected, result)
		}
	}
}

func TestPhase3_IfStatement(t *testing.T) {
	// 簡単な if 文のテスト
	input := "if 1 > 0 { x := 42 }"

	sc := scanner.NewScanner(input)
	parser := parser.NewParser(sc)
	stmt := parser.ParseStatement()

	env := eval.NewEnvironment()
	eval.EvalStatement(stmt, env)

	// x が設定されているかチェック
	value, exists := env.Get("x")
	if !exists {
		t.Errorf("variable x should exist")
	}
	if value != 42 {
		t.Errorf("expected x to be 42, got %d", value)
	}
}

func TestPhase3_IfElseStatement(t *testing.T) {
	// if-else 文のテスト
	input := "if 0 > 1 { x := 42 } else { x := 24 }"

	sc := scanner.NewScanner(input)
	parser := parser.NewParser(sc)
	stmt := parser.ParseStatement()

	env := eval.NewEnvironment()
	eval.EvalStatement(stmt, env)

	// x が else ブロックの値になっているかチェック
	value, exists := env.Get("x")
	if !exists {
		t.Errorf("variable x should exist")
	}
	if value != 24 {
		t.Errorf("expected x to be 24, got %d", value)
	}
}

func TestPhase3_SimpleVariableComparison(t *testing.T) {
	// 変数を使った比較のテスト
	env := eval.NewEnvironment()
	env.Set("x", 10)
	env.Set("y", 5)

	sc := scanner.NewScanner("x > y")
	parser := parser.NewParser(sc)
	ast := parser.ParseExpression()
	result := eval.EvalWithEnvironment(ast, env)

	if result != 1 {
		t.Errorf("expected 1 (true), got %d", result)
	}
}
