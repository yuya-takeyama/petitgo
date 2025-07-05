//go:build linux && amd64
// +build linux,amd64

package asmgen

import (
	"strings"
	"testing"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/token"
)

func TestX86_64Generator_Generate_SimpleMain(t *testing.T) {
	gen := NewX86_64Generator()

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

	// Check x86_64 assembly structure
	if !strings.Contains(result, ".section .text") {
		t.Error("Missing text section")
	}
	if !strings.Contains(result, ".globl _start") {
		t.Error("Missing global start symbol")
	}
	if !strings.Contains(result, "_start:") {
		t.Error("Missing start function label")
	}
	if !strings.Contains(result, "movq $42, %rax") {
		t.Error("Missing number loading")
	}
	if !strings.Contains(result, "call _print_number") {
		t.Error("Missing print_number call")
	}
}

func TestX86_64Generator_Generate_BinaryOperations(t *testing.T) {
	tests := []struct {
		name     string
		operator token.Token
		expected string
	}{
		{"addition", token.ADD, "addq %rbx, %rax"},
		{"subtraction", token.SUB, "subq %rbx, %rax"},
		{"multiplication", token.MUL, "imulq %rbx, %rax"},
		{"division", token.QUO, "idivq %rbx"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewX86_64Generator()

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

func TestX86_64Generator_GenerateRuntime(t *testing.T) {
	gen := NewX86_64Generator()
	runtime := gen.GenerateRuntime()

	// Check x86_64 runtime functions
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

	// Check x86_64 specific implementation details
	if !strings.Contains(runtime, "movb $48, (%rsi)") {
		t.Error("Missing ASCII '0' conversion")
	}
	if !strings.Contains(runtime, "idivq %rbx") {
		t.Error("Missing division operation")
	}
	if !strings.Contains(runtime, "syscall") {
		t.Error("Missing Linux system call")
	}
}

func TestX86_64Generator_Generate_IfStatement(t *testing.T) {
	gen := NewX86_64Generator()

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

	// Check x86_64 if statement generation
	if !strings.Contains(result, "cmpq %rbx, %rax") {
		t.Error("Missing comparison instruction")
	}
	if !strings.Contains(result, "setg %al") {
		t.Error("Missing condition evaluation")
	}
	if !strings.Contains(result, "testq %rax, %rax") {
		t.Error("Missing test instruction")
	}
	if !strings.Contains(result, "jz") {
		t.Error("Missing conditional jump")
	}
}
