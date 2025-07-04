package eval

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

func TestTypeSystem_IntegerArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1 + 2", 3},
		{"10 - 5", 5},
		{"3 * 4", 12},
		{"8 / 2", 4},
	}

	for _, tt := range tests {
		s := scanner.NewScanner(tt.input)
		p := parser.NewParser(s)
		expr := p.ParseExpression()

		result := EvalValue(expr)

		intVal, ok := result.(*IntValue)
		if !ok {
			t.Fatalf("Input: %s, expected IntValue, got %T", tt.input, result)
		}

		if intVal.Value != tt.expected {
			t.Errorf("Input: %s, expected: %d, got: %d", tt.input, tt.expected, intVal.Value)
		}
	}
}

func TestTypeSystem_StringConcatenation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello" + "world"`, "helloworld"},
		{`"Go" + "lang"`, "Golang"},
		{`"" + "test"`, "test"},
	}

	for _, tt := range tests {
		s := scanner.NewScanner(tt.input)
		p := parser.NewParser(s)
		expr := p.ParseExpression()

		result := EvalValue(expr)

		stringVal, ok := result.(*StringValue)
		if !ok {
			t.Fatalf("Input: %s, expected StringValue, got %T", tt.input, result)
		}

		if stringVal.Value != tt.expected {
			t.Errorf("Input: %s, expected: %s, got: %s", tt.input, tt.expected, stringVal.Value)
		}
	}
}

func TestTypeSystem_MixedTypeConcatenation(t *testing.T) {
	// When mixing types, they should be converted to strings
	input := `42 + "world"`
	s := scanner.NewScanner(input)
	p := parser.NewParser(s)
	expr := p.ParseExpression()

	result := EvalValue(expr)

	stringVal, ok := result.(*StringValue)
	if !ok {
		t.Fatalf("Expected StringValue for mixed type addition, got %T", result)
	}

	expected := "42world"
	if stringVal.Value != expected {
		t.Errorf("Expected: %s, got: %s", expected, stringVal.Value)
	}
}

func TestTypeSystem_BooleanComparisons(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"5 > 3", true},
		{"3 > 5", false},
		{"5 == 5", true},
		{"5 != 3", true},
		{"5 <= 5", true},
		{"3 >= 5", false},
	}

	for _, tt := range tests {
		s := scanner.NewScanner(tt.input)
		p := parser.NewParser(s)
		expr := p.ParseExpression()

		result := EvalValue(expr)

		boolVal, ok := result.(*BoolValue)
		if !ok {
			t.Fatalf("Input: %s, expected BoolValue, got %T", tt.input, result)
		}

		if boolVal.Value != tt.expected {
			t.Errorf("Input: %s, expected: %v, got: %v", tt.input, tt.expected, boolVal.Value)
		}
	}
}

func TestTypeSystem_StringComparisons(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"apple" == "apple"`, true},
		{`"apple" != "banana"`, true},
		{`"apple" < "banana"`, true},
		{`"banana" > "apple"`, true},
	}

	for _, tt := range tests {
		s := scanner.NewScanner(tt.input)
		p := parser.NewParser(s)
		expr := p.ParseExpression()

		result := EvalValue(expr)

		boolVal, ok := result.(*BoolValue)
		if !ok {
			t.Fatalf("Input: %s, expected BoolValue, got %T", tt.input, result)
		}

		if boolVal.Value != tt.expected {
			t.Errorf("Input: %s, expected: %v, got: %v", tt.input, tt.expected, boolVal.Value)
		}
	}
}

func TestTypeSystem_VariableWithTypes(t *testing.T) {
	env := NewEnvironment()

	// Set different types of variables
	env.Set("intVar", &IntValue{Value: 42})
	env.Set("stringVar", &StringValue{Value: "hello"})
	env.Set("boolVar", &BoolValue{Value: true})

	tests := []struct {
		varName  string
		expected string // expected type
	}{
		{"intVar", "int"},
		{"stringVar", "string"},
		{"boolVar", "bool"},
	}

	for _, tt := range tests {
		s := scanner.NewScanner(tt.varName)
		p := parser.NewParser(s)
		expr := p.ParseExpression()

		result := EvalValueWithEnvironment(expr, env)

		if result.Type() != tt.expected {
			t.Errorf("Variable %s, expected type: %s, got: %s", tt.varName, tt.expected, result.Type())
		}
	}
}

func TestTypeSystem_ArithmeticTypeErrors(t *testing.T) {
	// Test arithmetic operations with incompatible types
	env := NewEnvironment()
	env.Set("stringVar", &StringValue{Value: "hello"})
	env.Set("intVar", &IntValue{Value: 42})

	// string - int should return 0 (type error)
	input := "stringVar - intVar"
	s := scanner.NewScanner(input)
	p := parser.NewParser(s)
	expr := p.ParseExpression()

	result := EvalValueWithEnvironment(expr, env)

	intVal, ok := result.(*IntValue)
	if !ok {
		t.Fatalf("Expected IntValue for type error result, got %T", result)
	}

	if intVal.Value != 0 {
		t.Errorf("Expected 0 for type error, got %d", intVal.Value)
	}
}
