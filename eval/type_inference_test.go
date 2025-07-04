package eval

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

func TestTypeInference_VarStatement(t *testing.T) {
	tests := []struct {
		statement    string
		varName      string
		expectedType string
		expectedVal  interface{}
	}{
		{"var x int = 42", "x", "int", 42},
		{`var message string = "hello"`, "message", "string", "hello"},
		{"var flag bool = true", "flag", "bool", true},
		{"var result int = 10 + 20", "result", "int", 30},
	}

	for _, tt := range tests {
		env := NewEnvironment()

		s := scanner.NewScanner(tt.statement)
		p := parser.NewParser(s)
		stmt := p.ParseStatement()

		EvalStatement(stmt, env)

		// Check if variable was stored with correct type
		value, exists := env.Get(tt.varName)
		if !exists {
			t.Fatalf("Statement: %s, variable %s should exist", tt.statement, tt.varName)
		}

		if value.Type() != tt.expectedType {
			t.Errorf("Statement: %s, expected type %s, got %s", tt.statement, tt.expectedType, value.Type())
		}

		// Check value
		switch expectedVal := tt.expectedVal.(type) {
		case int:
			if intVal, ok := value.(*IntValue); ok {
				if intVal.Value != expectedVal {
					t.Errorf("Statement: %s, expected value %d, got %d", tt.statement, expectedVal, intVal.Value)
				}
			} else {
				t.Errorf("Statement: %s, expected IntValue, got %T", tt.statement, value)
			}
		case string:
			if strVal, ok := value.(*StringValue); ok {
				if strVal.Value != expectedVal {
					t.Errorf("Statement: %s, expected value %s, got %s", tt.statement, expectedVal, strVal.Value)
				}
			} else {
				t.Errorf("Statement: %s, expected StringValue, got %T", tt.statement, value)
			}
		case bool:
			if boolVal, ok := value.(*BoolValue); ok {
				if boolVal.Value != expectedVal {
					t.Errorf("Statement: %s, expected value %v, got %v", tt.statement, expectedVal, boolVal.Value)
				}
			} else {
				t.Errorf("Statement: %s, expected BoolValue, got %T", tt.statement, value)
			}
		}
	}
}

func TestTypeInference_AssignStatement(t *testing.T) {
	env := NewEnvironment()

	// First, assign an integer
	s := scanner.NewScanner("x := 42")
	p := parser.NewParser(s)
	stmt := p.ParseStatement()
	EvalStatement(stmt, env)

	value, exists := env.Get("x")
	if !exists {
		t.Fatal("Variable x should exist after assignment")
	}

	if value.Type() != "int" {
		t.Errorf("Expected type int, got %s", value.Type())
	}

	if intVal, ok := value.(*IntValue); ok {
		if intVal.Value != 42 {
			t.Errorf("Expected value 42, got %d", intVal.Value)
		}
	}

	// Now reassign with a string
	s = scanner.NewScanner(`x = "hello"`)
	p = parser.NewParser(s)
	stmt = p.ParseStatement()
	EvalStatement(stmt, env)

	value, exists = env.Get("x")
	if !exists {
		t.Fatal("Variable x should still exist after reassignment")
	}

	if value.Type() != "string" {
		t.Errorf("Expected type string after reassignment, got %s", value.Type())
	}

	if strVal, ok := value.(*StringValue); ok {
		if strVal.Value != "hello" {
			t.Errorf("Expected value 'hello', got %s", strVal.Value)
		}
	}
}

func TestTypeInference_IfStatementWithTruthy(t *testing.T) {
	tests := []struct {
		setup     string
		ifStmt    string
		varName   string
		shouldSet bool
	}{
		{`x := 1`, `if x { y := 42 }`, "y", true},       // non-zero int is truthy
		{`x := 0`, `if x { y := 42 }`, "y", false},      // zero int is falsy
		{`x := "hello"`, `if x { y := 42 }`, "y", true}, // non-empty string is truthy
		{`x := ""`, `if x { y := 42 }`, "y", false},     // empty string is falsy
		{`x := true`, `if x { y := 42 }`, "y", true},    // true is truthy
		{`x := false`, `if x { y := 42 }`, "y", false},  // false is falsy
	}

	for _, tt := range tests {
		env := NewEnvironment()

		// Setup
		s := scanner.NewScanner(tt.setup)
		p := parser.NewParser(s)
		stmt := p.ParseStatement()
		EvalStatement(stmt, env)

		// Execute if statement
		s = scanner.NewScanner(tt.ifStmt)
		p = parser.NewParser(s)
		stmt = p.ParseStatement()
		EvalStatement(stmt, env)

		// Check result
		_, exists := env.Get(tt.varName)
		if exists != tt.shouldSet {
			t.Errorf("Setup: %s, If: %s, expected var %s to exist: %v, but got: %v",
				tt.setup, tt.ifStmt, tt.varName, tt.shouldSet, exists)
		}
	}
}

func TestTypeInference_StringConcatenationInference(t *testing.T) {
	env := NewEnvironment()

	// Set variables with different types
	statements := []string{
		`name := "Alice"`,
		`age := 25`,
		`greeting := "Hello " + name`,         // string concatenation
		`info := name + " is " + "years old"`, // mixed concatenation
	}

	for _, stmt := range statements {
		s := scanner.NewScanner(stmt)
		p := parser.NewParser(s)
		statement := p.ParseStatement()
		EvalStatement(statement, env)
	}

	// Check types
	tests := []struct {
		varName      string
		expectedType string
	}{
		{"name", "string"},
		{"age", "int"},
		{"greeting", "string"},
		{"info", "string"},
	}

	for _, tt := range tests {
		value, exists := env.Get(tt.varName)
		if !exists {
			t.Fatalf("Variable %s should exist", tt.varName)
		}

		if value.Type() != tt.expectedType {
			t.Errorf("Variable %s: expected type %s, got %s", tt.varName, tt.expectedType, value.Type())
		}
	}

	// Check specific values
	if greeting, exists := env.Get("greeting"); exists {
		if strVal, ok := greeting.(*StringValue); ok {
			if strVal.Value != "Hello Alice" {
				t.Errorf("Expected greeting 'Hello Alice', got '%s'", strVal.Value)
			}
		}
	}
}
