package main

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/eval"
	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

// Test RED: This should fail initially - testing complete function workflow
func TestIntegration_SimpleFunctionDefinitionAndCall(t *testing.T) {
	env := eval.NewEnvironment()

	// Step 1: Define a simple function
	funcDefInput := "func double(x int) int { return x * 2 }"
	s1 := scanner.NewScanner(funcDefInput)
	p1 := parser.NewParser(s1)
	funcStmt := p1.ParseStatement()
	eval.EvalStatement(funcStmt, env)

	// Step 2: Call the function
	callInput := "double(5)"
	s2 := scanner.NewScanner(callInput)
	p2 := parser.NewParser(s2)
	callExpr := p2.ParseExpression()

	result := eval.EvalWithEnvironment(callExpr, env)

	// This should equal 10
	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}
}

// Test RED: This will definitely fail - testing more complex function
func TestIntegration_FunctionWithMultipleParameters(t *testing.T) {
	env := eval.NewEnvironment()

	// Define function with multiple parameters
	funcDefInput := "func add3(a int, b int, c int) int { return a + b + c }"
	s1 := scanner.NewScanner(funcDefInput)
	p1 := parser.NewParser(s1)
	funcStmt := p1.ParseStatement()
	eval.EvalStatement(funcStmt, env)

	// Call with 3 arguments
	callInput := "add3(1, 2, 3)"
	s2 := scanner.NewScanner(callInput)
	p2 := parser.NewParser(s2)
	callExpr := p2.ParseExpression()

	result := eval.EvalWithEnvironment(callExpr, env)

	// This should equal 6
	if result != 6 {
		t.Errorf("Expected 6, got %d", result)
	}
}

// Test RED: This will likely fail - nested function calls
func TestIntegration_NestedFunctionCalls(t *testing.T) {
	env := eval.NewEnvironment()

	// Define functions
	eval.EvalStatement(mustParseStatement("func double(x int) int { return x * 2 }"), env)
	eval.EvalStatement(mustParseStatement("func add(a int, b int) int { return a + b }"), env)

	// Call nested: add(double(3), double(4)) should be add(6, 8) = 14
	result := eval.EvalWithEnvironment(mustParseExpression("add(double(3), double(4))"), env)

	if result != 14 {
		t.Errorf("Expected 14, got %d", result)
	}
}

// Helper functions to reduce boilerplate
func mustParseStatement(input string) ast.Statement {
	s := scanner.NewScanner(input)
	p := parser.NewParser(s)
	return p.ParseStatement()
}

func mustParseExpression(input string) ast.ASTNode {
	s := scanner.NewScanner(input)
	p := parser.NewParser(s)
	return p.ParseExpression()
}
