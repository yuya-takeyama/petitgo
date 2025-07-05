package asmgen

import (
	"runtime"

	"github.com/yuya-takeyama/petitgo/ast"
)

// ArchGenerator defines the interface for architecture-specific assembly generators
type ArchGenerator interface {
	Generate(statements []ast.Statement) string
	GenerateRuntime() string
}

// AsmGenerator wraps the architecture-specific generator
type AsmGenerator struct {
	generator ArchGenerator
}

// NewAsmGenerator creates a new assembly generator for the current platform
func NewAsmGenerator() *AsmGenerator {
	var generator ArchGenerator

	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		generator = NewARM64Generator()
	} else if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		generator = NewX86_64Generator()
	} else {
		// Default to ARM64 for unsupported platforms
		generator = NewARM64Generator()
	}

	return &AsmGenerator{
		generator: generator,
	}
}

// Generate converts an AST to assembly code using the platform-specific generator
func (g *AsmGenerator) Generate(statements []ast.Statement) string {
	return g.generator.Generate(statements)
}

// GenerateRuntime generates helper functions using the platform-specific generator
func (g *AsmGenerator) GenerateRuntime() string {
	return g.generator.GenerateRuntime()
}
