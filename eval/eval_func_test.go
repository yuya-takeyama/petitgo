package eval

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

func TestEval_FunctionDefinitionAndCall(t *testing.T) {
	// Define function: func add(x int, y int) int { return x + y }
	funcDefInput := "func add(x int, y int) int { return x + y }"
	s1 := scanner.NewScanner(funcDefInput)
	p1 := parser.NewParser(s1)
	funcStmt := p1.ParseStatement()

	env := NewEnvironment()
	EvalStatement(funcStmt, env)

	// Call function: add(5, 3)
	callInput := "add(5, 3)"
	s2 := scanner.NewScanner(callInput)
	p2 := parser.NewParser(s2)
	callExpr := p2.ParseExpression()

	result := EvalWithEnvironment(callExpr, env)
	if result != 8 {
		t.Errorf("function call result wrong. expected=8, got=%d", result)
	}
}

func TestEval_FunctionWithoutReturn(t *testing.T) {
	// Define function without explicit return: func test() int { 42 }
	funcDefInput := "func test() int { 42 }"
	s1 := scanner.NewScanner(funcDefInput)
	p1 := parser.NewParser(s1)
	funcStmt := p1.ParseStatement()

	env := NewEnvironment()
	EvalStatement(funcStmt, env)

	// Call function: test()
	callInput := "test()"
	s2 := scanner.NewScanner(callInput)
	p2 := parser.NewParser(s2)
	callExpr := p2.ParseExpression()

	result := EvalWithEnvironment(callExpr, env)
	if result != 0 {
		t.Errorf("function call result wrong. expected=0, got=%d", result)
	}
}

func TestEval_ReturnStatement(t *testing.T) {
	// Test that return statement throws correct exception
	defer func() {
		if r := recover(); r != nil {
			if returnEx, ok := r.(*ReturnException); ok {
				if returnEx.Value != 42 {
					t.Errorf("return value wrong. expected=42, got=%d", returnEx.Value)
				}
			} else {
				t.Errorf("expected ReturnException, got %T", r)
			}
		} else {
			t.Error("expected panic for return statement")
		}
	}()

	returnInput := "return 42"
	s := scanner.NewScanner(returnInput)
	p := parser.NewParser(s)
	returnStmt := p.ParseStatement()

	env := NewEnvironment()
	EvalStatement(returnStmt, env)
}

func TestEval_FunctionRecursion(t *testing.T) {
	// Define recursive function: func fact(n int) int { if n <= 1 { return 1 } return n * fact(n - 1) }
	// For simplicity, we'll use a simpler recursive function
	funcDefInput := "func countdown(n int) int { if n <= 0 { return 0 } return countdown(n - 1) }"
	s1 := scanner.NewScanner(funcDefInput)
	p1 := parser.NewParser(s1)
	funcStmt := p1.ParseStatement()

	env := NewEnvironment()
	EvalStatement(funcStmt, env)

	// Call function: countdown(3)
	callInput := "countdown(3)"
	s2 := scanner.NewScanner(callInput)
	p2 := parser.NewParser(s2)
	callExpr := p2.ParseExpression()

	result := EvalWithEnvironment(callExpr, env)
	if result != 0 {
		t.Errorf("recursive function call result wrong. expected=0, got=%d", result)
	}
}
