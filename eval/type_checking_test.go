package eval

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

func TestTypeChecking_VarStatement_CorrectTypes(t *testing.T) {
	tests := []struct {
		input        string
		expectedType string
		expectedVal  interface{}
	}{
		{"var x int = 42", "int", 42},
		{"var s string = \"hello\"", "string", "hello"},
		{"var b bool = true", "bool", true},
	}

	for _, tt := range tests {
		env := NewEnvironment()

		// Parse and evaluate the var statement
		s := scanner.NewScanner(tt.input)
		p := parser.NewParser(s)
		stmt := p.ParseStatement()
		EvalStatement(stmt, env)

		// Check the stored value
		value, exists := env.Get("x")
		if tt.input[4] == 's' { // if variable name is 's'
			value, exists = env.Get("s")
		} else if tt.input[4] == 'b' { // if variable name is 'b'
			value, exists = env.Get("b")
		}

		if !exists {
			t.Fatalf("Variable not found for input: %s", tt.input)
		}

		if value.Type() != tt.expectedType {
			t.Errorf("Input: %s, expected type: %s, got: %s", tt.input, tt.expectedType, value.Type())
		}

		switch v := value.(type) {
		case *IntValue:
			if v.Value != tt.expectedVal {
				t.Errorf("Input: %s, expected value: %v, got: %v", tt.input, tt.expectedVal, v.Value)
			}
		case *StringValue:
			if v.Value != tt.expectedVal {
				t.Errorf("Input: %s, expected value: %v, got: %v", tt.input, tt.expectedVal, v.Value)
			}
		case *BoolValue:
			if v.Value != tt.expectedVal {
				t.Errorf("Input: %s, expected value: %v, got: %v", tt.input, tt.expectedVal, v.Value)
			}
		}
	}
}

func TestTypeChecking_VarStatement_TypeMismatch(t *testing.T) {
	tests := []struct {
		input             string
		expectedType      string
		expectedZeroValue interface{}
	}{
		{"var x int = \"hello\"", "int", 0}, // string assigned to int -> zero value
		{"var s string = 42", "string", ""}, // int assigned to string -> zero value
		{"var b bool = 123", "bool", false}, // int assigned to bool -> zero value
	}

	for _, tt := range tests {
		env := NewEnvironment()

		// Parse and evaluate the var statement
		s := scanner.NewScanner(tt.input)
		p := parser.NewParser(s)
		stmt := p.ParseStatement()
		EvalStatement(stmt, env)

		// Check the stored value
		varName := string(tt.input[4]) // x, s, or b
		value, exists := env.Get(varName)

		if !exists {
			t.Fatalf("Variable not found for input: %s", tt.input)
		}

		if value.Type() != tt.expectedType {
			t.Errorf("Input: %s, expected type: %s, got: %s", tt.input, tt.expectedType, value.Type())
		}

		// Check that we got the zero value
		switch v := value.(type) {
		case *IntValue:
			if v.Value != tt.expectedZeroValue {
				t.Errorf("Input: %s, expected zero value: %v, got: %v", tt.input, tt.expectedZeroValue, v.Value)
			}
		case *StringValue:
			if v.Value != tt.expectedZeroValue {
				t.Errorf("Input: %s, expected zero value: %v, got: %v", tt.input, tt.expectedZeroValue, v.Value)
			}
		case *BoolValue:
			if v.Value != tt.expectedZeroValue {
				t.Errorf("Input: %s, expected zero value: %v, got: %v", tt.input, tt.expectedZeroValue, v.Value)
			}
		}
	}
}

func TestTypeChecking_FunctionParameters(t *testing.T) {
	env := NewEnvironment()

	// Define function: func add(a int, b int) int { return a + b }
	funcInput := "func add(a int, b int) int { return a + b }"
	s1 := scanner.NewScanner(funcInput)
	p1 := parser.NewParser(s1)
	funcStmt := p1.ParseStatement()
	EvalStatement(funcStmt, env)

	tests := []struct {
		callInput    string
		expectedType string
		expectedVal  int
	}{
		{"add(5, 3)", "int", 8},                 // correct types
		{"add(\"hello\", \"world\")", "int", 0}, // wrong types -> zero values used
	}

	for _, tt := range tests {
		// Call function
		s2 := scanner.NewScanner(tt.callInput)
		p2 := parser.NewParser(s2)
		callExpr := p2.ParseExpression()

		result := EvalWithEnvironment(callExpr, env)
		if result != tt.expectedVal {
			t.Errorf("Call: %s, expected: %d, got: %d", tt.callInput, tt.expectedVal, result)
		}
	}
}

func TestTypeChecking_FunctionMissingArguments(t *testing.T) {
	env := NewEnvironment()

	// Define function: func greet(name string, age int) int { return age }
	funcInput := "func greet(name string, age int) int { return age }"
	s1 := scanner.NewScanner(funcInput)
	p1 := parser.NewParser(s1)
	funcStmt := p1.ParseStatement()
	EvalStatement(funcStmt, env)

	// Call with missing arguments -> should use zero values
	callInput := "greet()"
	s2 := scanner.NewScanner(callInput)
	p2 := parser.NewParser(s2)
	callExpr := p2.ParseExpression()

	result := EvalWithEnvironment(callExpr, env)
	// Should return 0 (zero value for int parameter 'age')
	if result != 0 {
		t.Errorf("Expected 0 for missing arguments, got: %d", result)
	}
}

func TestTypeInference_AssignStatements(t *testing.T) {
	tests := []struct {
		statements   []string
		varName      string
		expectedType string
		expectedVal  interface{}
	}{
		// Type inference from literal values
		{[]string{"x := 42"}, "x", "int", 42},
		{[]string{"name := \"Alice\""}, "name", "string", "Alice"},
		{[]string{"flag := true"}, "flag", "bool", true},

		// Type checking on reassignment (same type)
		{[]string{"x := 42", "x = 100"}, "x", "int", 100},
		{[]string{"name := \"Alice\"", "name = \"Bob\""}, "name", "string", "Bob"},

		// Type mismatch on reassignment (should use zero value of original type)
		{[]string{"x := 42", "x = \"hello\""}, "x", "int", 0},
		{[]string{"name := \"Alice\"", "name = 123"}, "name", "string", ""},
	}

	for _, tt := range tests {
		testEnv := NewEnvironment()

		// Execute all statements
		for _, stmt := range tt.statements {
			s := scanner.NewScanner(stmt)
			p := parser.NewParser(s)
			astStmt := p.ParseStatement()
			EvalStatement(astStmt, testEnv)
		}

		// Check the final value
		value, exists := testEnv.Get(tt.varName)
		if !exists {
			t.Fatalf("Variable %s not found after statements: %v", tt.varName, tt.statements)
		}

		if value.Type() != tt.expectedType {
			t.Errorf("Statements: %v, variable %s expected type: %s, got: %s",
				tt.statements, tt.varName, tt.expectedType, value.Type())
		}

		// Check value
		switch v := value.(type) {
		case *IntValue:
			if v.Value != tt.expectedVal {
				t.Errorf("Statements: %v, variable %s expected value: %v, got: %v",
					tt.statements, tt.varName, tt.expectedVal, v.Value)
			}
		case *StringValue:
			if v.Value != tt.expectedVal {
				t.Errorf("Statements: %v, variable %s expected value: %v, got: %v",
					tt.statements, tt.varName, tt.expectedVal, v.Value)
			}
		case *BoolValue:
			if v.Value != tt.expectedVal {
				t.Errorf("Statements: %v, variable %s expected value: %v, got: %v",
					tt.statements, tt.varName, tt.expectedVal, v.Value)
			}
		}
	}
}
