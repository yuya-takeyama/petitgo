package asmgen

import (
	"runtime"
	"strings"
	"testing"

	"github.com/yuya-takeyama/petitgo/ast"
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
