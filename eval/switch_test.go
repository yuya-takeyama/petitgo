package eval

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

func TestEvalSwitchStatement(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		setupX   int
		expected string
	}{
		{
			name: "switch with matching case",
			input: `switch x {
case 1:
	result = "one"
case 2:
	result = "two"
default:
	result = "other"
}`,
			setupX:   1,
			expected: "one",
		},
		{
			name: "switch with second case match",
			input: `switch x {
case 1:
	result = "one"
case 2:
	result = "two"
default:
	result = "other"
}`,
			setupX:   2,
			expected: "two",
		},
		{
			name: "switch with default case",
			input: `switch x {
case 1:
	result = "one"
case 2:
	result = "two"
default:
	result = "other"
}`,
			setupX:   99,
			expected: "other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()

			// Setup environment directly
			env.Set("x", &IntValue{Value: tt.setupX})
			env.Set("result", &StringValue{Value: ""})

			// Parse and evaluate switch statement
			sc := scanner.NewScanner(tt.input)
			p := parser.NewParser(sc)
			stmt := p.ParseStatement()

			if stmt == nil {
				t.Fatalf("expected statement, got nil")
			}

			EvalStatement(stmt, env)

			// Check result
			resultVal, exists := env.Get("result")
			if !exists {
				t.Fatalf("result variable not found in environment")
			}

			strVal, ok := resultVal.(*StringValue)
			if !ok {
				t.Fatalf("expected StringValue, got %T", resultVal)
			}

			if strVal.Value != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, strVal.Value)
			}
		})
	}
}

// Test switch with various types
func TestEvalSwitchTypes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		setup func(env *Environment)
		check func(t *testing.T, env *Environment)
	}{
		{
			name: "switch with string values",
			input: `switch name {
case "alice":
	result = "Hello Alice"
case "bob":
	result = "Hello Bob"
default:
	result = "Hello stranger"
}`,
			setup: func(env *Environment) {
				env.Set("name", &StringValue{Value: "alice"})
				env.Set("result", &StringValue{Value: ""})
			},
			check: func(t *testing.T, env *Environment) {
				resultVal, _ := env.Get("result")
				result := resultVal.(*StringValue)
				if result.Value != "Hello Alice" {
					t.Errorf("expected 'Hello Alice', got '%s'", result.Value)
				}
			},
		},
		{
			name: "switch with boolean",
			input: `switch isValid {
case true:
	result = "valid"
case false:
	result = "invalid"
}`,
			setup: func(env *Environment) {
				env.Set("isValid", &BoolValue{Value: true})
				env.Set("result", &StringValue{Value: ""})
			},
			check: func(t *testing.T, env *Environment) {
				resultVal, _ := env.Get("result")
				result := resultVal.(*StringValue)
				if result.Value != "valid" {
					t.Errorf("expected 'valid', got '%s'", result.Value)
				}
			},
		},
		{
			name: "switch with no match and no default",
			input: `switch x {
case 1:
	result = 10
case 2:
	result = 20
}`,
			setup: func(env *Environment) {
				env.Set("x", &IntValue{Value: 3})
				env.Set("result", &IntValue{Value: 99})
			},
			check: func(t *testing.T, env *Environment) {
				resultVal, _ := env.Get("result")
				result := resultVal.(*IntValue)
				if result.Value != 99 {
					t.Errorf("expected unchanged value 99, got %d", result.Value)
				}
			},
		},
		{
			name: "switch with multiple statements in case",
			input: `switch x {
case 1:
	a = 10
	b = 20
	result = a + b
case 2:
	result = 100
default:
	result = 0
}`,
			setup: func(env *Environment) {
				env.Set("x", &IntValue{Value: 1})
				env.Set("a", &IntValue{Value: 0})
				env.Set("b", &IntValue{Value: 0})
				env.Set("result", &IntValue{Value: 0})
			},
			check: func(t *testing.T, env *Environment) {
				resultVal, _ := env.Get("result")
				result := resultVal.(*IntValue)
				if result.Value != 30 {
					t.Errorf("expected 30, got %d", result.Value)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			tt.setup(env)

			sc := scanner.NewScanner(tt.input)
			p := parser.NewParser(sc)
			stmt := p.ParseStatement()

			if stmt == nil {
				t.Fatalf("expected statement, got nil")
			}

			EvalStatement(stmt, env)
			tt.check(t, env)
		})
	}
}
