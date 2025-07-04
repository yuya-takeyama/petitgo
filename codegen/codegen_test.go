package codegen

import (
	"strings"
	"testing"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/token"
)

func TestGenerateNode(t *testing.T) {
	tests := []struct {
		name     string
		node     ast.ASTNode
		expected string
	}{
		{
			name:     "number node",
			node:     &ast.NumberNode{Value: 42},
			expected: "42",
		},
		{
			name:     "string node",
			node:     &ast.StringNode{Value: "hello"},
			expected: `"hello"`,
		},
		{
			name:     "boolean true",
			node:     &ast.BooleanNode{Value: true},
			expected: "true",
		},
		{
			name:     "boolean false",
			node:     &ast.BooleanNode{Value: false},
			expected: "false",
		},
		{
			name:     "variable node",
			node:     &ast.VariableNode{Name: "x"},
			expected: "x",
		},
		{
			name: "binary operation - addition",
			node: &ast.BinaryOpNode{
				Left:     &ast.NumberNode{Value: 1},
				Operator: token.ADD,
				Right:    &ast.NumberNode{Value: 2},
			},
			expected: "(1 + 2)",
		},
		{
			name: "binary operation - nested",
			node: &ast.BinaryOpNode{
				Left: &ast.BinaryOpNode{
					Left:     &ast.NumberNode{Value: 2},
					Operator: token.ADD,
					Right:    &ast.NumberNode{Value: 3},
				},
				Operator: token.MUL,
				Right:    &ast.NumberNode{Value: 4},
			},
			expected: "((2 + 3) * 4)",
		},
		{
			name: "function call - print",
			node: &ast.CallNode{
				Function: "print",
				Arguments: []ast.ASTNode{
					&ast.NumberNode{Value: 123},
				},
			},
			expected: "println(123)",
		},
		{
			name: "function call - custom",
			node: &ast.CallNode{
				Function: "add",
				Arguments: []ast.ASTNode{
					&ast.NumberNode{Value: 1},
					&ast.NumberNode{Value: 2},
				},
			},
			expected: "add(1, 2)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGenerator()
			result := g.Generate(tt.node)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGenerateStatement(t *testing.T) {
	tests := []struct {
		name     string
		stmt     ast.Statement
		expected string
	}{
		{
			name: "var statement",
			stmt: &ast.VarStatement{
				Name:     "x",
				TypeName: "int",
				Value:    &ast.NumberNode{Value: 42},
			},
			expected: "var x int = 42",
		},
		{
			name: "assign statement",
			stmt: &ast.AssignStatement{
				Name:  "x",
				Value: &ast.NumberNode{Value: 10},
			},
			expected: "x := 10",
		},
		{
			name: "expression statement",
			stmt: &ast.ExpressionStatement{
				Expression: &ast.CallNode{
					Function: "print",
					Arguments: []ast.ASTNode{
						&ast.StringNode{Value: "hello"},
					},
				},
			},
			expected: `println("hello")`,
		},
		{
			name: "return statement with value",
			stmt: &ast.ReturnStatement{
				Value: &ast.NumberNode{Value: 42},
			},
			expected: "return 42",
		},
		{
			name:     "return statement without value",
			stmt:     &ast.ReturnStatement{},
			expected: "return",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGenerator()
			// Use a single statement wrapped in a program
			code := g.GenerateProgram([]ast.Statement{tt.stmt})

			// Extract the generated statement (skip package/main wrapper)
			lines := strings.Split(code, "\n")
			var actualLine string
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed != "" &&
					!strings.HasPrefix(trimmed, "package") &&
					!strings.HasPrefix(trimmed, "func main") &&
					trimmed != "}" {
					actualLine = trimmed
					break
				}
			}

			if actualLine != tt.expected {
				t.Errorf("Expected %q, got %q\nFull output:\n%s", tt.expected, actualLine, code)
			}
		})
	}
}

func TestGenerateProgram(t *testing.T) {
	tests := []struct {
		name       string
		statements []ast.Statement
		contains   []string
	}{
		{
			name: "basic program",
			statements: []ast.Statement{
				&ast.AssignStatement{
					Name:  "x",
					Value: &ast.NumberNode{Value: 10},
				},
				&ast.ExpressionStatement{
					Expression: &ast.CallNode{
						Function: "print",
						Arguments: []ast.ASTNode{
							&ast.VariableNode{Name: "x"},
						},
					},
				},
			},
			contains: []string{
				"package main",
				"func main()",
				"x := 10",
				"println(x)",
			},
		},
		{
			name: "program with main function",
			statements: []ast.Statement{
				&ast.FuncStatement{
					Name: "main",
					Body: &ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.CallNode{
									Function: "print",
									Arguments: []ast.ASTNode{
										&ast.StringNode{Value: "Hello"},
									},
								},
							},
						},
					},
				},
			},
			contains: []string{
				"package main",
				"func main() {",
				`println("Hello")`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGenerator()
			result := g.GenerateProgram(tt.statements)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected output to contain %q\nFull output:\n%s", expected, result)
				}
			}
		})
	}
}

func TestParenthesesPreservation(t *testing.T) {
	tests := []struct {
		name     string
		node     ast.ASTNode
		expected string
	}{
		{
			name: "simple parentheses",
			node: &ast.BinaryOpNode{
				Left: &ast.BinaryOpNode{
					Left:     &ast.NumberNode{Value: 2},
					Operator: token.ADD,
					Right:    &ast.NumberNode{Value: 3},
				},
				Operator: token.MUL,
				Right:    &ast.NumberNode{Value: 4},
			},
			expected: "((2 + 3) * 4)",
		},
		{
			name: "complex precedence",
			node: &ast.BinaryOpNode{
				Left:     &ast.NumberNode{Value: 10},
				Operator: token.ADD,
				Right: &ast.BinaryOpNode{
					Left:     &ast.NumberNode{Value: 5},
					Operator: token.MUL,
					Right:    &ast.NumberNode{Value: 2},
				},
			},
			expected: "(10 + (5 * 2))",
		},
		{
			name: "deeply nested",
			node: &ast.BinaryOpNode{
				Left: &ast.BinaryOpNode{
					Left: &ast.BinaryOpNode{
						Left:     &ast.NumberNode{Value: 1},
						Operator: token.ADD,
						Right:    &ast.NumberNode{Value: 2},
					},
					Operator: token.MUL,
					Right:    &ast.NumberNode{Value: 3},
				},
				Operator: token.ADD,
				Right:    &ast.NumberNode{Value: 4},
			},
			expected: "(((1 + 2) * 3) + 4)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGenerator()
			result := g.Generate(tt.node)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
