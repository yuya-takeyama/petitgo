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

	// Generate all functions first
	for _, stmt := range statements {
		if funcStmt, ok := stmt.(*ast.FuncStatement); ok {
			g.generateFunction(funcStmt)
		}
	}

	return g.output.String()
}

func (g *AsmGenerator) generateFunction(funcStmt *ast.FuncStatement) {
	// Reset variables for each function
	g.variables = make(map[string]int)
	g.stackSize = 0

	g.writeLine(fmt.Sprintf("_%s:", funcStmt.Name))
	g.writeLine("    // Function prologue")
	g.writeLine("    stp x29, x30, [sp, #-16]!") // Save frame pointer and link register
	g.writeLine("    mov x29, sp")               // Set frame pointer
	g.writeLine("    sub sp, sp, #64")           // Reserve 64 bytes for local variables

	// For functions with parameters, handle the first parameter in x0
	if len(funcStmt.Parameters) > 0 {
		// Store parameter in stack
		g.stackSize += 8
		g.variables[funcStmt.Parameters[0].Name] = g.stackSize
		g.writeLine(fmt.Sprintf("    // Parameter: %s", funcStmt.Parameters[0].Name))
		g.writeLine(fmt.Sprintf("    str x0, [x29, #-%d]", g.stackSize))
	}

	// Generate function body
	g.generateBlock(funcStmt.Body)

	// Function epilogue only for main (other functions use explicit return)
	if funcStmt.Name == "main" {
		g.writeLine("    // Exit")
		g.writeLine("    add sp, sp, #64")         // Restore stack space
		g.writeLine("    ldp x29, x30, [sp], #16") // Restore frame pointer and link register
		g.writeLine("    mov x0, #0")              // exit status
		g.writeLine("    mov x16, #1")             // sys_exit
		g.writeLine("    svc #0x80")               // system call
	}
	g.writeLine("")
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
		// Check if variable already exists
		if offset, exists := g.variables[s.Name]; exists {
			// Variable reassignment
			g.writeLine(fmt.Sprintf("    // %s = value (reassignment)", s.Name))
			g.generateExpression(s.Value)
			g.writeLine(fmt.Sprintf("    str x0, [x29, #-%d]", offset))
		} else {
			// New variable allocation
			g.stackSize += 8
			g.variables[s.Name] = g.stackSize

			g.writeLine(fmt.Sprintf("    // %s := value", s.Name))
			g.generateExpression(s.Value)
			offset := g.stackSize
			g.writeLine(fmt.Sprintf("    str x0, [x29, #-%d]", offset)) // store relative to frame pointer
		}
	case *ast.ReassignStatement:
		// Variable reassignment
		if offset, exists := g.variables[s.Name]; exists {
			g.writeLine(fmt.Sprintf("    // %s = value", s.Name))
			g.generateExpression(s.Value)
			g.writeLine(fmt.Sprintf("    str x0, [x29, #-%d]", offset))
		}
	case *ast.IfStatement:
		g.generateIfStatement(s)
	case *ast.ForStatement:
		g.generateForStatement(s)
	case *ast.ReturnStatement:
		g.generateReturnStatement(s)
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
	case *ast.CallNode:
		g.generateFunctionCall(n)
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
		case token.LSS:
			g.writeLine("    cmp x1, x0")
			g.writeLine("    cset x0, lo") // set x0 to 1 if x1 < x0
		case token.GTR:
			g.writeLine("    cmp x1, x0")
			g.writeLine("    cset x0, hi") // set x0 to 1 if x1 > x0
		case token.LEQ:
			g.writeLine("    cmp x1, x0")
			g.writeLine("    cset x0, ls") // set x0 to 1 if x1 <= x0
		case token.GEQ:
			g.writeLine("    cmp x1, x0")
			g.writeLine("    cset x0, hs") // set x0 to 1 if x1 >= x0
		case token.EQL:
			g.writeLine("    cmp x1, x0")
			g.writeLine("    cset x0, eq") // set x0 to 1 if x1 == x0
		case token.NEQ:
			g.writeLine("    cmp x1, x0")
			g.writeLine("    cset x0, ne") // set x0 to 1 if x1 != x0
		}
	}
}

func (g *AsmGenerator) generateIfStatement(stmt *ast.IfStatement) {
	elseLabel := g.getNewLabel()
	endLabel := g.getNewLabel()

	// Generate condition
	g.writeLine("    // if condition")
	g.generateExpression(stmt.Condition)
	g.writeLine("    cmp x0, #0")                     // compare with 0
	g.writeLine(fmt.Sprintf("    beq %s", elseLabel)) // branch if equal (false)

	// Generate then block
	g.writeLine("    // then block")
	g.generateBlock(stmt.ThenBlock)
	g.writeLine(fmt.Sprintf("    b %s", endLabel)) // branch to end

	// Generate else block (if exists)
	g.writeLine(fmt.Sprintf("%s:", elseLabel))
	if stmt.ElseBlock != nil {
		g.writeLine("    // else block")
		g.generateBlock(stmt.ElseBlock)
	}

	g.writeLine(fmt.Sprintf("%s:", endLabel))
}

func (g *AsmGenerator) generateForStatement(stmt *ast.ForStatement) {
	loopLabel := g.getNewLabel()
	endLabel := g.getNewLabel()

	g.writeLine(fmt.Sprintf("%s:", loopLabel))

	// Generate condition
	g.generateExpression(stmt.Condition)
	g.writeLine("    cmp x0, #0")
	g.writeLine(fmt.Sprintf("    beq %s", endLabel)) // exit if false

	// Generate body
	g.generateBlock(stmt.Body)

	g.writeLine(fmt.Sprintf("    b %s", loopLabel)) // loop back
	g.writeLine(fmt.Sprintf("%s:", endLabel))
}

func (g *AsmGenerator) generateReturnStatement(stmt *ast.ReturnStatement) {
	if stmt.Value != nil {
		g.generateExpression(stmt.Value)
	}
	// Return to caller (result already in x0)
	g.writeLine("    add sp, sp, #64")
	g.writeLine("    ldp x29, x30, [sp], #16")
	g.writeLine("    ret")
}

func (g *AsmGenerator) generateFunctionCall(call *ast.CallNode) {
	// For now, only handle recursive function calls
	// Save caller's registers
	g.writeLine("    // Function call: " + call.Function)

	// Generate argument (assume single argument for fibonacci)
	if len(call.Arguments) > 0 {
		g.generateExpression(call.Arguments[0])
	}

	// Call function
	g.writeLine(fmt.Sprintf("    bl _%s", call.Function))
}

func (g *AsmGenerator) getNewLabel() string {
	g.labelNum++
	return fmt.Sprintf("L%d", g.labelNum)
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
