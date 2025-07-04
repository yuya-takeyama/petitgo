package main

import "testing"

func TestParseVarStatement(t *testing.T) {
	input := "var x int = 42"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	stmt := parser.ParseStatement()

	varStmt, ok := stmt.(*VarStatement)
	if !ok {
		t.Fatalf("expected VarStatement, got %T", stmt)
	}

	if varStmt.Name != "x" {
		t.Errorf("expected variable name 'x', got %s", varStmt.Name)
	}

	if varStmt.TypeName != "int" {
		t.Errorf("expected type 'int', got %s", varStmt.TypeName)
	}

	numberNode, ok := varStmt.Value.(*NumberNode)
	if !ok {
		t.Fatalf("expected NumberNode, got %T", varStmt.Value)
	}

	if numberNode.Value != 42 {
		t.Errorf("expected value 42, got %d", numberNode.Value)
	}
}

func TestParseAssignStatement(t *testing.T) {
	input := "x := 42"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	stmt := parser.ParseStatement()

	assignStmt, ok := stmt.(*AssignStatement)
	if !ok {
		t.Fatalf("expected AssignStatement, got %T", stmt)
	}

	if assignStmt.Name != "x" {
		t.Errorf("expected variable name 'x', got %s", assignStmt.Name)
	}

	numberNode, ok := assignStmt.Value.(*NumberNode)
	if !ok {
		t.Fatalf("expected NumberNode, got %T", assignStmt.Value)
	}

	if numberNode.Value != 42 {
		t.Errorf("expected value 42, got %d", numberNode.Value)
	}
}

func TestParseSimpleAssignment(t *testing.T) {
	input := "x = 10"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	stmt := parser.ParseStatement()

	assignStmt, ok := stmt.(*AssignStatement)
	if !ok {
		t.Fatalf("expected AssignStatement, got %T", stmt)
	}

	if assignStmt.Name != "x" {
		t.Errorf("expected variable name 'x', got %s", assignStmt.Name)
	}

	numberNode, ok := assignStmt.Value.(*NumberNode)
	if !ok {
		t.Fatalf("expected NumberNode, got %T", assignStmt.Value)
	}

	if numberNode.Value != 10 {
		t.Errorf("expected value 10, got %d", numberNode.Value)
	}
}

func TestParseVariableReference(t *testing.T) {
	input := "x"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	expr := parser.ParseExpression()

	varRef, ok := expr.(*VariableNode)
	if !ok {
		t.Fatalf("expected VariableNode, got %T", expr)
	}

	if varRef.Name != "x" {
		t.Errorf("expected variable name 'x', got %s", varRef.Name)
	}
}
