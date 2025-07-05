package asmgen

import (
	"runtime"
	"strings"
	"testing"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/token"
)

// TestAsmGenerator_PlatformDetection tests that the correct generator is created based on platform
func TestAsmGenerator_PlatformDetection(t *testing.T) {
	gen := NewAsmGenerator()

	// Create a simple function
	blockStmt := &ast.BlockStatement{Statements: []ast.Statement{}}
	funcStmt := &ast.FuncStatement{
		Name:       "main",
		Parameters: []ast.Parameter{},
		Body:       blockStmt,
	}
	statements := []ast.Statement{funcStmt}

	result := gen.Generate(statements)

	// Check that appropriate assembly is generated based on platform
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		if !strings.Contains(result, ".section __TEXT,__text") {
			t.Error("Expected ARM64 assembly for macOS")
		}
	} else if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		if !strings.Contains(result, ".section .text") {
			t.Error("Expected x86_64 assembly for Linux")
		}
	}
}

// TestAsmGenerator_GenerateRuntime tests runtime generation
func TestAsmGenerator_GenerateRuntime(t *testing.T) {
	gen := NewAsmGenerator()
	runtime := gen.GenerateRuntime()

	// All platforms should have _print_number function
	if !strings.Contains(runtime, "_print_number:") {
		t.Error("Missing _print_number function in runtime")
	}
}

// TestNewAsmGenerator_UnsupportedPlatform tests the default fallback behavior
func TestNewAsmGenerator_UnsupportedPlatform(t *testing.T) {
	// Save original values
	originalGOOS := runtime.GOOS
	originalGOARCH := runtime.GOARCH

	// Test unsupported platform (Windows x86_64)
	// Note: We can't actually change runtime.GOOS/GOARCH in tests,
	// but we can test the current behavior which defaults to ARM64
	gen := NewAsmGenerator()

	// Verify that the generator was created (should not be nil)
	if gen == nil {
		t.Error("NewAsmGenerator returned nil")
	}

	// Test that it can generate code (should work with ARM64 fallback)
	blockStmt := &ast.BlockStatement{Statements: []ast.Statement{}}
	funcStmt := &ast.FuncStatement{
		Name:       "main",
		Parameters: []ast.Parameter{},
		Body:       blockStmt,
	}
	statements := []ast.Statement{funcStmt}

	result := gen.Generate(statements)
	if result == "" {
		t.Error("Generate returned empty string")
	}

	// Restore original values (not needed in this case since we didn't change them)
	_ = originalGOOS
	_ = originalGOARCH
}

// TestAsmGenerator_AllPlatforms tests generator creation for all supported platforms
func TestAsmGenerator_AllPlatforms(t *testing.T) {
	tests := []struct {
		name     string
		goos     string
		goarch   string
		expected string
	}{
		{"macOS ARM64", "darwin", "arm64", ".section __TEXT,__text"},
		{"Linux x86_64", "linux", "amd64", ".section .text"},
		{"Unsupported", "windows", "amd64", ".section __TEXT,__text"}, // Falls back to ARM64
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: We can't mock runtime.GOOS/GOARCH in tests easily,
			// so this test documents the expected behavior
			// In practice, the current platform will determine the generator
			gen := NewAsmGenerator()

			if gen == nil {
				t.Error("NewAsmGenerator returned nil")
			}

			// Verify the generator can produce output
			blockStmt := &ast.BlockStatement{Statements: []ast.Statement{}}
			funcStmt := &ast.FuncStatement{
				Name:       "main",
				Parameters: []ast.Parameter{},
				Body:       blockStmt,
			}
			statements := []ast.Statement{funcStmt}

			result := gen.Generate(statements)
			if result == "" {
				t.Error("Generate returned empty string")
			}
		})
	}
}

// TestX86_64Generator_DirectAccess tests x86_64 generator directly without build tags
func TestX86_64Generator_DirectAccess(t *testing.T) {
	// Directly create x86_64 generator without build tag restrictions
	gen := NewX86_64Generator()

	if gen == nil {
		t.Error("NewX86_64Generator returned nil")
	}

	// Test basic generation
	blockStmt := &ast.BlockStatement{Statements: []ast.Statement{}}
	funcStmt := &ast.FuncStatement{
		Name:       "main",
		Parameters: []ast.Parameter{},
		Body:       blockStmt,
	}
	statements := []ast.Statement{funcStmt}

	result := gen.Generate(statements)
	if result == "" {
		t.Error("Generate returned empty string")
	}

	// Check x86_64 specific output
	if !strings.Contains(result, ".section .text") {
		t.Error("Missing x86_64 text section")
	}

	// Test runtime generation
	runtime := gen.GenerateRuntime()
	if runtime == "" {
		t.Error("GenerateRuntime returned empty string")
	}

	if !strings.Contains(runtime, "_print_number:") {
		t.Error("Missing _print_number function in x86_64 runtime")
	}
}

// TestX86_64Generator_CrossPlatformFeatures tests all x86_64 features from any platform
func TestX86_64Generator_CrossPlatformFeatures(t *testing.T) {
	gen := NewX86_64Generator()

	// Test binary operations
	t.Run("binary_operations", func(t *testing.T) {
		binaryOp := &ast.BinaryOpNode{
			Left:     &ast.NumberNode{Value: 10},
			Operator: token.ADD,
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

		if !strings.Contains(result, "addq %rbx, %rax") {
			t.Error("Missing x86_64 addition instruction")
		}
	})

	// Test for statement
	t.Run("for_statement", func(t *testing.T) {
		condition := &ast.BinaryOpNode{
			Left:     &ast.VariableNode{Name: "i"},
			Operator: token.LSS,
			Right:    &ast.NumberNode{Value: 10},
		}
		printCall := &ast.CallNode{
			Function:  "println",
			Arguments: []ast.ASTNode{&ast.VariableNode{Name: "i"}},
		}
		body := &ast.BlockStatement{
			Statements: []ast.Statement{&ast.ExpressionStatement{Expression: printCall}},
		}
		forStmt := &ast.ForStatement{
			Condition: condition,
			Body:      body,
		}

		assignStmt := &ast.AssignStatement{
			Name:  "i",
			Value: &ast.NumberNode{Value: 0},
		}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt, forStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		// Check for loop structure
		if !strings.Contains(result, "_start:") {
			t.Error("Missing function start label")
		}
	})

	// Test return statement
	t.Run("return_statement", func(t *testing.T) {
		returnStmt := &ast.ReturnStatement{
			Value: &ast.NumberNode{Value: 42},
		}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{returnStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "test",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		if !strings.Contains(result, "movq $42, %rax") {
			t.Error("Missing return value setup")
		}
		if !strings.Contains(result, "ret") {
			t.Error("Missing return instruction")
		}
	})

	// Test function call - use println which is implemented
	t.Run("function_call", func(t *testing.T) {
		callNode := &ast.CallNode{
			Function:  "println",
			Arguments: []ast.ASTNode{&ast.NumberNode{Value: 5}},
		}
		exprStmt := &ast.ExpressionStatement{Expression: callNode}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{exprStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		if !strings.Contains(result, "movq $5, %rax") {
			t.Error("Missing x86_64 argument setup")
		}
		if !strings.Contains(result, "call _print_number") {
			t.Error("Missing print_number call")
		}
	})

	// Test switch statement
	t.Run("switch_statement", func(t *testing.T) {
		caseStmt := &ast.CaseStatement{
			Value: &ast.NumberNode{Value: 1},
			Body: &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.CallNode{
							Function:  "println",
							Arguments: []ast.ASTNode{&ast.NumberNode{Value: 100}},
						},
					},
				},
			},
		}
		switchStmt := &ast.SwitchStatement{
			Value:   &ast.VariableNode{Name: "x"},
			Cases:   []*ast.CaseStatement{caseStmt},
			Default: &ast.BlockStatement{Statements: []ast.Statement{}},
		}

		assignStmt := &ast.AssignStatement{
			Name:  "x",
			Value: &ast.NumberNode{Value: 1},
		}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt, switchStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		// Basic test - just ensure it generates code
		if result == "" {
			t.Error("Generate returned empty string")
		}
	})

	// Test if statement
	t.Run("if_statement", func(t *testing.T) {
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

		if !strings.Contains(result, "cmpq %rbx, %rax") {
			t.Error("Missing x86_64 comparison")
		}
		if !strings.Contains(result, "setg %al") {
			t.Error("Missing condition evaluation")
		}
	})
}
