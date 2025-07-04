package asmgen

import (
	"fmt"
	"strings"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/token"
)

// AsmGenerator generates ARM64 assembly code for M1 Mac
type AsmGenerator struct {
	output    strings.Builder
	labelNum  int
	variables map[string]int // variable name -> stack offset
	stackSize int
}

// NewAsmGenerator creates a new assembly generator
func NewAsmGenerator() *AsmGenerator {
	return &AsmGenerator{
		variables: make(map[string]int),
		stackSize: 0,
	}
}

// Generate converts an AST to ARM64 assembly code
func (g *AsmGenerator) Generate(statements []ast.Statement) string {
	g.output.Reset()

	// Basic Mach-O sections for macOS ARM64
	g.writeLine(".section __TEXT,__text,regular,pure_instructions")
	g.writeLine(".globl _main")
	g.writeLine(".p2align 2")
	g.writeLine("")
	g.writeLine("_main:")
	g.writeLine("    // Proper ARM64 stack frame setup")
	g.writeLine("    stp x29, x30, [sp, #-16]!") // Save frame pointer and link register
	g.writeLine("    mov x29, sp")               // Set frame pointer
	g.writeLine("    sub sp, sp, #64")           // Reserve 64 bytes for local variables

	// Generate code for statements
	for _, stmt := range statements {
		if funcStmt, ok := stmt.(*ast.FuncStatement); ok && funcStmt.Name == "main" {
			g.generateBlock(funcStmt.Body)
		}
	}

	// Exit - restore stack frame and exit
	g.writeLine("    // Exit")
	g.writeLine("    add sp, sp, #64")         // Restore stack space
	g.writeLine("    ldp x29, x30, [sp], #16") // Restore frame pointer and link register
	g.writeLine("    mov x0, #0")              // exit status
	g.writeLine("    mov x16, #1")             // sys_exit
	g.writeLine("    svc #0x80")               // system call

	return g.output.String()
}

func (g *AsmGenerator) generateBlock(block *ast.BlockStatement) {
	for _, stmt := range block.Statements {
		g.generateStatement(stmt)
	}
}

func (g *AsmGenerator) generateStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		if callNode, ok := s.Expression.(*ast.CallNode); ok {
			if callNode.Function == "println" && len(callNode.Arguments) > 0 {
				g.generatePrintln(callNode.Arguments[0])
			}
		}
	case *ast.AssignStatement:
		// Allocate stack space for variable (frame-pointer based)
		g.stackSize += 8
		g.variables[s.Name] = g.stackSize

		g.writeLine(fmt.Sprintf("    // %s := value", s.Name))
		g.generateExpression(s.Value)
		offset := g.stackSize
		g.writeLine(fmt.Sprintf("    str x0, [x29, #-%d]", offset)) // store relative to frame pointer
	}
}

func (g *AsmGenerator) generatePrintln(arg ast.ASTNode) {
	// Generate expression and get result in x0
	g.generateExpression(arg)

	// Convert x0 to character and print (simplified for single digits)
	g.writeLine("    // Convert number to ASCII and print")
	g.writeLine("    add x0, x0, #48") // Add '0' to convert to ASCII
	g.writeLine("    bl _print_char")
	g.writeLine("    mov x0, #10") // newline
	g.writeLine("    bl _print_char")
}

func (g *AsmGenerator) generateExpression(expr ast.ASTNode) {
	switch n := expr.(type) {
	case *ast.NumberNode:
		g.writeLine(fmt.Sprintf("    mov x0, #%d", n.Value))
	case *ast.VariableNode:
		if offset, exists := g.variables[n.Name]; exists {
			g.writeLine(fmt.Sprintf("    ldr x0, [x29, #-%d]", offset)) // load relative to frame pointer
		}
	case *ast.BinaryOpNode:
		// Generate left operand
		g.generateExpression(n.Left)
		g.writeLine("    str x0, [sp, #-16]!") // push result with proper alignment

		// Generate right operand
		g.generateExpression(n.Right)
		g.writeLine("    ldr x1, [sp], #16") // pop left result to x1 with proper alignment

		// Perform operation
		switch n.Operator {
		case token.ADD:
			g.writeLine("    add x0, x1, x0")
		case token.SUB:
			g.writeLine("    sub x0, x1, x0")
		case token.MUL:
			g.writeLine("    mul x0, x1, x0")
		}
	}
}

func (g *AsmGenerator) writeLine(s string) {
	g.output.WriteString(s + "\n")
}

// GenerateRuntime generates helper functions for ARM64 macOS
func (g *AsmGenerator) GenerateRuntime() string {
	runtime := `
_print_char:
    // Print character in x0 to stdout (macOS ARM64)
    // Save character to stack
    str x0, [sp, #-16]!
    
    // write syscall
    mov x0, #1              // stdout
    mov x1, sp              // buffer (character on stack)
    mov x2, #1              // count
    mov x16, #4             // sys_write
    svc #0x80
    
    // Restore stack
    add sp, sp, #16
    ret
`
	return runtime
}
