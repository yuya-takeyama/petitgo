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
	case *ast.SwitchStatement:
		g.generateSwitchStatement(s)
	}
}

func (g *AsmGenerator) generatePrintln(arg ast.ASTNode) {
	// Generate expression and get result in x0
	g.generateExpression(arg)

	// Call print_number function to handle multi-digit numbers
	g.writeLine("    // Print number")
	g.writeLine("    bl _print_number")
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
		case token.QUO:
			g.writeLine("    udiv x0, x1, x0")
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
	case *ast.FieldAccessNode:
		g.generateFieldAccess(n)
	case *ast.SliceLiteral:
		g.generateSliceLiteral(n)
	case *ast.IndexAccess:
		g.generateIndexAccess(n)
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

func (g *AsmGenerator) generateSwitchStatement(stmt *ast.SwitchStatement) {
	endLabel := g.getNewLabel()

	// Generate the switch value and store it
	g.writeLine("    // switch expression")
	g.generateExpression(stmt.Value)
	g.writeLine("    str x0, [sp, #-16]!") // Store switch value on stack

	// Generate case comparisons and labels
	caseLabels := make([]string, len(stmt.Cases))
	for i := range stmt.Cases {
		caseLabels[i] = g.getNewLabel()
	}

	defaultLabel := g.getNewLabel()
	if stmt.Default != nil {
		defaultLabel = g.getNewLabel()
	}

	// Compare each case value
	for i, caseStmt := range stmt.Cases {
		g.writeLine(fmt.Sprintf("    // case %d comparison", i))
		g.generateExpression(caseStmt.Value)
		g.writeLine("    ldr x1, [sp]") // Load switch value from stack
		g.writeLine("    cmp x1, x0")
		g.writeLine(fmt.Sprintf("    beq %s", caseLabels[i]))
	}

	// If no case matched, go to default or end
	if stmt.Default != nil {
		g.writeLine(fmt.Sprintf("    b %s", defaultLabel))
	} else {
		g.writeLine(fmt.Sprintf("    b %s", endLabel))
	}

	// Generate case bodies
	for i, caseStmt := range stmt.Cases {
		g.writeLine(fmt.Sprintf("%s:", caseLabels[i]))
		g.writeLine(fmt.Sprintf("    // case %d body", i))
		g.generateBlock(caseStmt.Body)
		g.writeLine(fmt.Sprintf("    b %s", endLabel)) // break to end
	}

	// Generate default case if exists
	if stmt.Default != nil {
		g.writeLine(fmt.Sprintf("%s:", defaultLabel))
		g.writeLine("    // default case")
		g.generateBlock(stmt.Default)
	}

	// End label
	g.writeLine(fmt.Sprintf("%s:", endLabel))
	g.writeLine("    add sp, sp, #16") // Restore stack (remove stored switch value)
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

func (g *AsmGenerator) generateFieldAccess(node *ast.FieldAccessNode) {
	// For ARM64 assembly, we'll need to implement struct field access
	// For now, we'll store structs as consecutive memory locations
	// and access fields by offset

	g.writeLine("    // Field access: obj.field")

	// Generate the object expression first
	g.generateExpression(node.Object)

	// For now, simple field access - assume field offset is fixed
	// This is a simplified implementation
	// In a real implementation, we'd need struct type information
	fieldOffset := getFieldOffset(node.Field)
	g.writeLine(fmt.Sprintf("    // Access field '%s' at offset %d", node.Field, fieldOffset))
	g.writeLine(fmt.Sprintf("    ldr x0, [x0, #%d]", fieldOffset))
}

// getFieldOffset returns a simple field offset (simplified implementation)
func getFieldOffset(fieldName string) int {
	// For simplicity, we'll use a fixed offset based on field name
	// In a real implementation, this would come from struct type info
	switch fieldName {
	case "x", "name", "first":
		return 0
	case "y", "age", "second":
		return 8
	case "z", "third":
		return 16
	default:
		return 0
	}
}

func (g *AsmGenerator) generateSliceLiteral(node *ast.SliceLiteral) {
	// For ARM64 assembly, we'll implement slices as dynamic arrays
	// For now, we'll allocate static memory and store element count
	g.writeLine("    // Slice literal creation")

	// For simplicity, return the length of the slice
	// In a real implementation, we'd allocate heap memory
	elementCount := len(node.Elements)
	g.writeLine(fmt.Sprintf("    mov x0, #%d", elementCount))
}

func (g *AsmGenerator) generateIndexAccess(node *ast.IndexAccess) {
	// For ARM64 assembly, we'll access array elements by offset
	g.writeLine("    // Index access: arr[index]")

	// Generate the object (array/slice) expression
	g.generateExpression(node.Object)
	g.writeLine("    str x0, [sp, #-16]!") // Store array base address

	// Generate the index expression
	g.generateExpression(node.Index)
	g.writeLine("    ldr x1, [sp], #16") // Load array base address

	// For simplicity, assume 8-byte elements (int64)
	// Calculate offset: index * 8
	g.writeLine("    lsl x0, x0, #3") // x0 = index * 8
	g.writeLine("    add x0, x1, x0") // x0 = base + offset
	g.writeLine("    ldr x0, [x0]")   // Load value at address
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
_print_number:
    // Print number in x0 as decimal (ARM64 macOS) - Simple recursive approach
    stp x29, x30, [sp, #-16]!
    mov x29, sp
    
    // Handle zero case
    cmp x0, #0
    bne print_nonzero
    mov x0, #48             // ASCII '0'
    bl _print_char
    b print_done
    
print_nonzero:
    // If number >= 10, print higher digits first (recursion)
    cmp x0, #10
    blo print_single_digit
    
    // Save current number
    str x0, [sp, #-16]!
    
    // Divide by 10 and recurse
    mov x1, #10
    udiv x0, x0, x1
    bl _print_number
    
    // Restore number and print last digit
    ldr x0, [sp], #16
    mov x1, #10
    udiv x2, x0, x1         // x2 = x0 / 10
    msub x0, x2, x1, x0     // x0 = x0 - (x2 * 10) = remainder
    
print_single_digit:
    // Convert single digit to ASCII and print
    add x0, x0, #48         // Convert to ASCII
    bl _print_char
    
print_done:
    ldp x29, x30, [sp], #16
    ret

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
