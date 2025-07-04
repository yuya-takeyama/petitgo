package parser

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/scanner"
)

func TestParser_FuncStatement(t *testing.T) {
	input := "func add(x int, y int) int { return x + y }"
	s := scanner.NewScanner(input)
	p := NewParser(s)

	stmt := p.ParseStatement()
	funcStmt, ok := stmt.(*ast.FuncStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.FuncStatement. got=%T", stmt)
	}

	if funcStmt.Name != "add" {
		t.Errorf("function name wrong. expected=add, got=%s", funcStmt.Name)
	}

	if len(funcStmt.Parameters) != 2 {
		t.Errorf("parameters count wrong. expected=2, got=%d", len(funcStmt.Parameters))
	}

	if funcStmt.Parameters[0].Name != "x" || funcStmt.Parameters[0].Type != "int" {
		t.Errorf("first parameter wrong. expected=x int, got=%s %s",
			funcStmt.Parameters[0].Name, funcStmt.Parameters[0].Type)
	}

	if funcStmt.Parameters[1].Name != "y" || funcStmt.Parameters[1].Type != "int" {
		t.Errorf("second parameter wrong. expected=y int, got=%s %s",
			funcStmt.Parameters[1].Name, funcStmt.Parameters[1].Type)
	}

	if funcStmt.ReturnType != "int" {
		t.Errorf("return type wrong. expected=int, got=%s", funcStmt.ReturnType)
	}

	if funcStmt.Body == nil {
		t.Error("function body is nil")
	}
}

func TestParser_ReturnStatement(t *testing.T) {
	input := "return 42"
	s := scanner.NewScanner(input)
	p := NewParser(s)

	stmt := p.ParseStatement()
	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ReturnStatement. got=%T", stmt)
	}

	if returnStmt.Value == nil {
		t.Error("return value is nil")
	}

	// Check if return value is a NumberNode with value 42
	if numberNode, ok := returnStmt.Value.(*ast.NumberNode); ok {
		if numberNode.Value != 42 {
			t.Errorf("return value wrong. expected=42, got=%d", numberNode.Value)
		}
	} else {
		t.Errorf("return value is not NumberNode. got=%T", returnStmt.Value)
	}
}

func TestParser_EmptyReturn(t *testing.T) {
	input := "return"
	s := scanner.NewScanner(input)
	p := NewParser(s)

	stmt := p.ParseStatement()
	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ReturnStatement. got=%T", stmt)
	}

	if returnStmt.Value != nil {
		t.Error("empty return should have nil value")
	}
}
