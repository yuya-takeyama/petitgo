package eval

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

func TestValue_StringEvaluation(t *testing.T) {
	input := `"hello world"`
	s := scanner.NewScanner(input)
	p := parser.NewParser(s)
	expr := p.ParseExpression()

	result := EvalValue(expr)

	// Should be StringValue
	stringVal, ok := result.(*StringValue)
	if !ok {
		t.Fatalf("Expected StringValue, got %T", result)
	}

	if stringVal.Value != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", stringVal.Value)
	}

	if stringVal.Type() != "string" {
		t.Errorf("Expected type 'string', got '%s'", stringVal.Type())
	}

	if stringVal.String() != "hello world" {
		t.Errorf("Expected String() 'hello world', got '%s'", stringVal.String())
	}
}

func TestValue_BooleanEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
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

		if boolVal.Type() != "bool" {
			t.Errorf("Expected type 'bool', got '%s'", boolVal.Type())
		}
	}
}

func TestValue_IntegerEvaluation(t *testing.T) {
	input := "42"
	s := scanner.NewScanner(input)
	p := parser.NewParser(s)
	expr := p.ParseExpression()

	result := EvalValue(expr)

	intVal, ok := result.(*IntValue)
	if !ok {
		t.Fatalf("Expected IntValue, got %T", result)
	}

	if intVal.Value != 42 {
		t.Errorf("Expected 42, got %d", intVal.Value)
	}

	if intVal.Type() != "int" {
		t.Errorf("Expected type 'int', got '%s'", intVal.Type())
	}

	if intVal.String() != "42" {
		t.Errorf("Expected String() '42', got '%s'", intVal.String())
	}
}
