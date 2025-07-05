package asmgen

import (
	"strings"
	"testing"

	"github.com/yuya-takeyama/petitgo/ast"
)

// Test switch statement code generation for ARM64
func TestGenerateSwitchStatementARM64(t *testing.T) {
	tests := []struct {
		name      string
		stmt      *ast.SwitchStatement
		expected  []string
		notExpect []string
	}{
		{
			name: "simple switch with two cases",
			stmt: &ast.SwitchStatement{
				Value: &ast.VariableNode{Name: "x"},
				Cases: []*ast.CaseStatement{
					{
						Value: &ast.NumberNode{Value: 1},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReassignStatement{
									Name:  "result",
									Value: &ast.NumberNode{Value: 10},
								},
							},
						},
					},
					{
						Value: &ast.NumberNode{Value: 2},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReassignStatement{
									Name:  "result",
									Value: &ast.NumberNode{Value: 20},
								},
							},
						},
					},
				},
			},
			expected: []string{
				"// switch expression",
				"// case 0 comparison",
				"mov x0, #1",
				"cmp x1, x0",
				"beq",
				"// case 1 comparison",
				"mov x0, #2",
				"// case 0 body",
				"mov x0, #10",
				"// case 1 body",
				"mov x0, #20",
			},
		},
		{
			name: "switch with default case",
			stmt: &ast.SwitchStatement{
				Value: &ast.VariableNode{Name: "x"},
				Cases: []*ast.CaseStatement{
					{
						Value: &ast.NumberNode{Value: 1},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReassignStatement{
									Name:  "result",
									Value: &ast.NumberNode{Value: 10},
								},
							},
						},
					},
				},
				Default: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.ReassignStatement{
							Name:  "result",
							Value: &ast.NumberNode{Value: 99},
						},
					},
				},
			},
			expected: []string{
				"// switch expression",
				"// case 0 comparison",
				"// default case",
				"mov x0, #99",
			},
		},
		{
			name: "switch with string cases",
			stmt: &ast.SwitchStatement{
				Value: &ast.VariableNode{Name: "name"},
				Cases: []*ast.CaseStatement{
					{
						Value: &ast.StringNode{Value: "alice"},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReassignStatement{
									Name:  "result",
									Value: &ast.NumberNode{Value: 1},
								},
							},
						},
					},
				},
			},
			expected: []string{
				"// switch expression",
				"// case 0 comparison",
				"str_0", // String label reference
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewARM64Generator()
			gen.variables = map[string]int{
				"x":      24,
				"name":   32,
				"result": 40,
			}

			gen.generateSwitchStatement(tt.stmt)
			output := gen.output.String()

			// Check expected strings
			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain '%s', but it didn't\nOutput:\n%s", expected, output)
				}
			}

			// Check strings that should not appear
			for _, notExpect := range tt.notExpect {
				if strings.Contains(output, notExpect) {
					t.Errorf("expected output NOT to contain '%s', but it did\nOutput:\n%s", notExpect, output)
				}
			}
		})
	}
}

// Test switch statement code generation for x86_64
func TestGenerateSwitchStatementX86_64(t *testing.T) {
	tests := []struct {
		name      string
		stmt      *ast.SwitchStatement
		expected  []string
		notExpect []string
	}{
		{
			name: "simple switch with two cases",
			stmt: &ast.SwitchStatement{
				Value: &ast.VariableNode{Name: "x"},
				Cases: []*ast.CaseStatement{
					{
						Value: &ast.NumberNode{Value: 1},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReassignStatement{
									Name:  "result",
									Value: &ast.NumberNode{Value: 10},
								},
							},
						},
					},
					{
						Value: &ast.NumberNode{Value: 2},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReassignStatement{
									Name:  "result",
									Value: &ast.NumberNode{Value: 20},
								},
							},
						},
					},
				},
			},
			expected: []string{
				"# switch expression",
				"# case 0 comparison",
				"movq $1, %rax",
				"cmpq %rax, %rbx",
				"je",
				"# case 1 comparison",
				"movq $2, %rax",
				"# case 0 body",
				"movq $10, %rax",
				"# case 1 body",
				"movq $20, %rax",
			},
		},
		{
			name: "switch with default case",
			stmt: &ast.SwitchStatement{
				Value: &ast.VariableNode{Name: "x"},
				Cases: []*ast.CaseStatement{
					{
						Value: &ast.NumberNode{Value: 1},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReassignStatement{
									Name:  "result",
									Value: &ast.NumberNode{Value: 10},
								},
							},
						},
					},
				},
				Default: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.ReassignStatement{
							Name:  "result",
							Value: &ast.NumberNode{Value: 99},
						},
					},
				},
			},
			expected: []string{
				"# switch expression",
				"# case 0 comparison",
				"# default case",
				"movq $99, %rax",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewX86_64Generator()
			gen.variables = map[string]int{
				"x":      24,
				"result": 40,
			}

			gen.generateSwitchStatement(tt.stmt)
			output := gen.output.String()

			// Check expected strings
			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain '%s', but it didn't\nOutput:\n%s", expected, output)
				}
			}

			// Check strings that should not appear
			for _, notExpect := range tt.notExpect {
				if strings.Contains(output, notExpect) {
					t.Errorf("expected output NOT to contain '%s', but it did\nOutput:\n%s", notExpect, output)
				}
			}
		})
	}
}

// Test switch edge cases
func TestGenerateSwitchEdgeCases(t *testing.T) {
	t.Run("switch with empty cases", func(t *testing.T) {
		stmt := &ast.SwitchStatement{
			Value: &ast.VariableNode{Name: "x"},
			Cases: []*ast.CaseStatement{},
			Default: &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.ReassignStatement{
						Name:  "result",
						Value: &ast.NumberNode{Value: 0},
					},
				},
			},
		}

		// ARM64
		genARM := NewARM64Generator()
		genARM.variables = map[string]int{"x": 24, "result": 32}
		genARM.generateSwitchStatement(stmt)
		outputARM := genARM.output.String()

		if !strings.Contains(outputARM, "// default case") {
			t.Errorf("ARM64: expected default case in output")
		}

		// x86_64
		genX86 := NewX86_64Generator()
		genX86.variables = map[string]int{"x": 24, "result": 32}
		genX86.generateSwitchStatement(stmt)
		outputX86 := genX86.output.String()

		if !strings.Contains(outputX86, "# default case") {
			t.Errorf("x86_64: expected default case in output")
		}
	})

	t.Run("switch with no default and no matching case", func(t *testing.T) {
		stmt := &ast.SwitchStatement{
			Value: &ast.VariableNode{Name: "x"},
			Cases: []*ast.CaseStatement{
				{
					Value: &ast.NumberNode{Value: 1},
					Body: &ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ReassignStatement{
								Name:  "result",
								Value: &ast.NumberNode{Value: 10},
							},
						},
					},
				},
			},
			Default: nil,
		}

		// ARM64
		genARM := NewARM64Generator()
		genARM.variables = map[string]int{"x": 24, "result": 32}
		genARM.generateSwitchStatement(stmt)
		outputARM := genARM.output.String()

		// Should jump to end label if no match
		if !strings.Contains(outputARM, "b L1") {
			t.Errorf("ARM64: expected branch to end label, got:\n%s", outputARM)
		}

		// x86_64
		genX86 := NewX86_64Generator()
		genX86.variables = map[string]int{"x": 24, "result": 32}
		genX86.generateSwitchStatement(stmt)
		outputX86 := genX86.output.String()

		// Should jump to end label if no match
		if !strings.Contains(outputX86, "jmp L1") {
			t.Errorf("x86_64: expected jump to end label, got:\n%s", outputX86)
		}
	})
}
