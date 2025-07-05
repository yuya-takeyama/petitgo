//go:build darwin && arm64
// +build darwin,arm64

package asmgen

import (
	"strings"
	"testing"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/token"
)

func TestARM64Generator_Generate_SimpleMain(t *testing.T) {
	gen := NewARM64Generator()

	// Create a simple main function with print statement
	printCall := &ast.CallNode{
		Function:  "println",
		Arguments: []ast.ASTNode{&ast.NumberNode{Value: 42}},
	}
	exprStmt := &ast.ExpressionStatement{Expression: printCall}
	blockStmt := &ast.BlockStatement{Statements: []ast.Statement{exprStmt}}

	funcStmt := &ast.FuncStatement{
		Name:       "main",
		Parameters: []ast.Parameter{},
		Body:       blockStmt,
	}

	statements := []ast.Statement{funcStmt}
	result := gen.Generate(statements)

	// Check ARM64 assembly structure
	if !strings.Contains(result, ".section __TEXT,__text,regular,pure_instructions") {
		t.Error("Missing TEXT section")
	}
	if !strings.Contains(result, ".globl _main") {
		t.Error("Missing global main symbol")
	}
	if !strings.Contains(result, "_main:") {
		t.Error("Missing main function label")
	}
	if !strings.Contains(result, "mov x0, #42") {
		t.Error("Missing number loading")
	}
	if !strings.Contains(result, "bl _print_number") {
		t.Error("Missing print_number call")
	}
}

func TestARM64Generator_Generate_BinaryOperations(t *testing.T) {
	tests := []struct {
		name     string
		operator token.Token
		expected string
	}{
		{"addition", token.ADD, "add x0, x0, x1"},
		{"subtraction", token.SUB, "sub x0, x0, x1"},
		{"multiplication", token.MUL, "mul x0, x0, x1"},
		{"division", token.QUO, "udiv x0, x0, x1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewARM64Generator()

			// Create: result := 10 + 5
			binaryOp := &ast.BinaryOpNode{
				Left:     &ast.NumberNode{Value: 10},
				Operator: tt.operator,
				Right:    &ast.NumberNode{Value: 5},
			}
			assignStmt := &ast.AssignStatement{
				Name:  "result",
				Value: binaryOp,
			}

			blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt}}
			funcStmt := &ast.FuncStatement{
				Name:       "main",
				Parameters: []ast.Parameter{},
				Body:       blockStmt,
			}

			statements := []ast.Statement{funcStmt}
			result := gen.Generate(statements)

			if !strings.Contains(result, tt.expected) {
				t.Errorf("Missing operation instruction: %s", tt.expected)
			}
		})
	}
}

func TestARM64Generator_GenerateRuntime(t *testing.T) {
	gen := NewARM64Generator()
	runtime := gen.GenerateRuntime()

	// Check ARM64 runtime functions
	expectedFunctions := []string{
		"_print_number:",
		"convert_loop:",
		"print_digits:",
		"print_loop:",
		"print_newline:",
	}

	for _, fn := range expectedFunctions {
		if !strings.Contains(runtime, fn) {
			t.Errorf("Missing runtime function: %s", fn)
		}
	}

	// Check ARM64 specific implementation details
	if !strings.Contains(runtime, "mov w3, #48") {
		t.Error("Missing ASCII '0' conversion")
	}
	if !strings.Contains(runtime, "udiv x4, x0, x3") {
		t.Error("Missing division operation")
	}
	if !strings.Contains(runtime, "msub x5, x4, x3, x0") {
		t.Error("Missing modulo operation")
	}
	if !strings.Contains(runtime, "svc #0x80") {
		t.Error("Missing macOS system call")
	}
}

func TestARM64Generator_Generate_IfStatement(t *testing.T) {
	gen := NewARM64Generator()

	// Create: if x > 5 { println(x) }
	condition := &ast.BinaryOpNode{
		Left:     &ast.VariableNode{Name: "x"},
		Operator: token.GTR,
		Right:    &ast.NumberNode{Value: 5},
	}

	printCall := &ast.CallNode{
		Function:  "println",
		Arguments: []ast.ASTNode{&ast.VariableNode{Name: "x"}},
	}
	thenBlock := &ast.BlockStatement{
		Statements: []ast.Statement{&ast.ExpressionStatement{Expression: printCall}},
	}

	ifStmt := &ast.IfStatement{
		Condition: condition,
		ThenBlock: thenBlock,
	}

	// Create variable x := 10
	assignStmt := &ast.AssignStatement{
		Name:  "x",
		Value: &ast.NumberNode{Value: 10},
	}

	blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt, ifStmt}}
	funcStmt := &ast.FuncStatement{
		Name:       "main",
		Parameters: []ast.Parameter{},
		Body:       blockStmt,
	}

	statements := []ast.Statement{funcStmt}
	result := gen.Generate(statements)

	// Check ARM64 if statement generation
	if !strings.Contains(result, "cmp x0, x1") {
		t.Error("Missing comparison instruction")
	}
	if !strings.Contains(result, "cset x0, gt") {
		t.Error("Missing condition evaluation")
	}
	if !strings.Contains(result, "cbz x0,") {
		t.Error("Missing conditional branch")
	}
}
