package asmgen

import (
	"fmt"
	"strings"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/token"
)

// ARM64Generator generates ARM64 assembly code for M1 Mac
type ARM64Generator struct {
	output         strings.Builder
	labelNum       int
	variables      map[string]int // variable name -> stack offset
	stackSize      int
	stringLiterals map[string]string // string value -> label name
	stringCount    int
}

// NewARM64Generator creates a new ARM64 assembly generator
func NewARM64Generator() *ARM64Generator {
	return &ARM64Generator{
		variables:      make(map[string]int),
		stackSize:      0,
		stringLiterals: make(map[string]string),
		stringCount:    0,
	}
}

// Generate converts an AST to ARM64 assembly code
func (g *ARM64Generator) Generate(statements []ast.Statement) string {
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

	// Generate string literals in data section
	g.generateStringLiterals()

	return g.output.String()
}

func (g *ARM64Generator) generateFunction(funcStmt *ast.FuncStatement) {
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

func (g *ARM64Generator) generateBlock(block *ast.BlockStatement) {
	for _, stmt := range block.Statements {
		g.generateStatement(stmt)
	}
}

func (g *ARM64Generator) generateStatement(stmt ast.Statement) {
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
			// New variable assignment
			g.stackSize += 8
			g.variables[s.Name] = g.stackSize
			g.writeLine(fmt.Sprintf("    // %s := value", s.Name))
			g.generateExpression(s.Value)
			g.writeLine(fmt.Sprintf("    str x0, [x29, #-%d]", g.stackSize))
		}
	case *ast.VarStatement:
		g.stackSize += 8
		g.variables[s.Name] = g.stackSize
		g.writeLine(fmt.Sprintf("    // var %s", s.Name))
		g.generateExpression(s.Value)
		g.writeLine(fmt.Sprintf("    str x0, [x29, #-%d]", g.stackSize))
	case *ast.IfStatement:
		g.generateIfStatement(s)
	case *ast.ForStatement:
		g.generateForStatement(s)
	case *ast.ReturnStatement:
		g.generateReturnStatement(s)
	case *ast.ReassignStatement:
		// Variable reassignment
		if offset, exists := g.variables[s.Name]; exists {
			g.writeLine(fmt.Sprintf("    // %s = value", s.Name))
			g.generateExpression(s.Value)
			g.writeLine(fmt.Sprintf("    str x0, [x29, #-%d]", offset))
		}
	case *ast.SwitchStatement:
		g.generateSwitchStatement(s)
	}
}

func (g *ARM64Generator) generatePrintln(arg ast.ASTNode) {
	g.generateExpression(arg)

	// Check argument type to determine print function
	switch arg.(type) {
	case *ast.StringNode:
		g.writeLine("    // Print string in x0")
		g.writeLine("    bl _print_string")
	default:
		g.writeLine("    // Print number in x0")
		g.writeLine("    bl _print_number")
	}
}

func (g *ARM64Generator) generateExpression(expr ast.ASTNode) {
	switch e := expr.(type) {
	case *ast.NumberNode:
		g.writeLine(fmt.Sprintf("    mov x0, #%d", e.Value))
	case *ast.VariableNode:
		if offset, exists := g.variables[e.Name]; exists {
			g.writeLine(fmt.Sprintf("    ldr x0, [x29, #-%d]", offset))
		}
	case *ast.StringNode:
		// Get or create string label
		label := g.getStringLabel(e.Value)
		g.writeLine(fmt.Sprintf("    adrp x0, %s@PAGE", label))
		g.writeLine(fmt.Sprintf("    add x0, x0, %s@PAGEOFF", label))
	case *ast.BinaryOpNode:
		// Left operand
		g.generateExpression(e.Left)
		g.writeLine("    str x0, [sp, #-16]!")

		// Right operand
		g.generateExpression(e.Right)
		g.writeLine("    mov x1, x0")
		g.writeLine("    ldr x0, [sp], #16")

		// Operation
		switch e.Operator {
		case token.ADD:
			g.writeLine("    add x0, x0, x1")
		case token.SUB:
			g.writeLine("    sub x0, x0, x1")
		case token.MUL:
			g.writeLine("    mul x0, x0, x1")
		case token.QUO:
			g.writeLine("    udiv x0, x0, x1")
		case token.EQL:
			g.writeLine("    cmp x0, x1")
			g.writeLine("    cset x0, eq")
		case token.NEQ:
			g.writeLine("    cmp x0, x1")
			g.writeLine("    cset x0, ne")
		case token.LSS:
			g.writeLine("    cmp x0, x1")
			g.writeLine("    cset x0, lt")
		case token.LEQ:
			g.writeLine("    cmp x0, x1")
			g.writeLine("    cset x0, le")
		case token.GTR:
			g.writeLine("    cmp x0, x1")
			g.writeLine("    cset x0, gt")
		case token.GEQ:
			g.writeLine("    cmp x0, x1")
			g.writeLine("    cset x0, ge")
		}
	case *ast.CallNode:
		g.generateFunctionCall(e)
	case *ast.FieldAccessNode:
		g.generateFieldAccess(e)
	case *ast.SliceLiteral:
		g.generateSliceLiteral(e)
	case *ast.IndexAccess:
		g.generateIndexAccess(e)
	}
}

func (g *ARM64Generator) generateIfStatement(stmt *ast.IfStatement) {
	endLabel := g.getNewLabel()

	g.generateExpression(stmt.Condition)
	g.writeLine("    cbz x0, " + endLabel)

	g.generateBlock(stmt.ThenBlock)

	if stmt.ElseBlock != nil {
		elseLabel := g.getNewLabel()
		g.writeLine("    b " + elseLabel)
		g.writeLine(endLabel + ":")
		g.generateBlock(stmt.ElseBlock)
		g.writeLine(elseLabel + ":")
	} else {
		g.writeLine(endLabel + ":")
	}
}

func (g *ARM64Generator) generateForStatement(stmt *ast.ForStatement) {
	startLabel := g.getNewLabel()
	endLabel := g.getNewLabel()

	g.writeLine(startLabel + ":")
	g.generateExpression(stmt.Condition)
	g.writeLine("    cbz x0, " + endLabel)

	g.generateBlock(stmt.Body)

	g.writeLine("    b " + startLabel)
	g.writeLine(endLabel + ":")
}

func (g *ARM64Generator) generateReturnStatement(stmt *ast.ReturnStatement) {
	if stmt.Value != nil {
		g.generateExpression(stmt.Value)
	}
	g.writeLine("    add sp, sp, #64")         // Restore stack space
	g.writeLine("    ldp x29, x30, [sp], #16") // Restore frame pointer and link register
	g.writeLine("    ret")                     // Return
}

func (g *ARM64Generator) generateSwitchStatement(stmt *ast.SwitchStatement) {
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

func (g *ARM64Generator) generateFunctionCall(call *ast.CallNode) {
	if len(call.Arguments) > 0 {
		g.generateExpression(call.Arguments[0])
	}
	g.writeLine(fmt.Sprintf("    bl _%s", call.Function))
}

func (g *ARM64Generator) generateFieldAccess(node *ast.FieldAccessNode) {
	g.writeLine("    // Field access: obj.field")
	g.generateExpression(node.Object)
	fieldOffset := getFieldOffset(node.Field)
	g.writeLine(fmt.Sprintf("    // Access field '%s' at offset %d", node.Field, fieldOffset))
	g.writeLine(fmt.Sprintf("    ldr x0, [x0, #%d]", fieldOffset))
}

func (g *ARM64Generator) generateSliceLiteral(node *ast.SliceLiteral) {
	g.writeLine("    // Slice literal creation")
	elementCount := len(node.Elements)
	g.writeLine(fmt.Sprintf("    mov x0, #%d", elementCount))
}

func (g *ARM64Generator) generateIndexAccess(node *ast.IndexAccess) {
	g.writeLine("    // Index access: arr[index]")
	g.generateExpression(node.Object)
	g.writeLine("    str x0, [sp, #-16]!") // Store array base address
	g.generateExpression(node.Index)
	g.writeLine("    ldr x1, [sp], #16") // Load array base address
	g.writeLine("    lsl x0, x0, #3")    // x0 = index * 8
	g.writeLine("    add x0, x1, x0")    // x0 = base + offset
	g.writeLine("    ldr x0, [x0]")      // Load value at address
}

func (g *ARM64Generator) getNewLabel() string {
	g.labelNum++
	return fmt.Sprintf("L%d", g.labelNum)
}

func (g *ARM64Generator) getStringLabel(value string) string {
	if label, exists := g.stringLiterals[value]; exists {
		return label
	}

	label := fmt.Sprintf("str_%d", g.stringCount)
	g.stringLiterals[value] = label
	g.stringCount++
	return label
}

func (g *ARM64Generator) generateStringLiterals() {
	if len(g.stringLiterals) == 0 {
		return
	}

	g.writeLine("")
	g.writeLine(".section __TEXT,__cstring,cstring_literals")

	for value, label := range g.stringLiterals {
		g.writeLine(fmt.Sprintf("%s:", label))
		g.writeLine(fmt.Sprintf("    .asciz \"%s\"", value))
	}
}

func (g *ARM64Generator) writeLine(s string) {
	g.output.WriteString(s + "\n")
}

// getFieldOffset returns a simple field offset (simplified implementation)
func getFieldOffset(fieldName string) int {
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

func (g *ARM64Generator) GenerateRuntime() string {
	runtime := `
// Runtime function to print numbers
.p2align 2
_print_number:
    stp x29, x30, [sp, #-16]!
    mov x29, sp
    sub sp, sp, #32
    
    // Save the number
    str x0, [x29, #-8]
    
    // Buffer for digits (16 bytes should be enough)
    add x1, x29, #-24  // Buffer pointer
    mov x2, #0         // Digit count
    
    // Handle zero case
    cmp x0, #0
    bne convert_loop
    mov w3, #48        // ASCII '0'
    strb w3, [x1]
    mov x2, #1
    b print_digits
    
.p2align 2
convert_loop:
    cbz x0, print_digits
    
    // Divide by 10
    mov x3, #10
    udiv x4, x0, x3    // x4 = x0 / 10
    msub x5, x4, x3, x0 // x5 = x0 - (x4 * 10) = x0 % 10
    
    // Convert digit to ASCII and store
    add w5, w5, #48    // Convert to ASCII
    strb w5, [x1, x2]  // Store digit
    add x2, x2, #1     // Increment digit count
    
    mov x0, x4         // Continue with quotient
    b convert_loop
    
.p2align 2
print_digits:
    // Print digits in reverse order
    cbz x2, print_newline
    
.p2align 2
print_loop:
    sub x2, x2, #1
    ldrb w3, [x1, x2]  // Load character to w3
    
    // Save registers before syscall
    stp x1, x2, [sp, #-16]!
    
    // Write system call
    sub sp, sp, #16    // Allocate space for character
    strb w3, [sp]      // Store character on stack
    mov x0, #1         // stdout
    mov x1, sp         // buffer (character on stack)
    mov x2, #1         // length
    mov x16, #4        // sys_write
    svc #0x80
    
    add sp, sp, #16    // Deallocate character space
    ldp x1, x2, [sp], #16  // Restore registers
    
    cbnz x2, print_loop
    
.p2align 2
print_newline:
    // Write newline character directly
    mov x16, #4        // sys_write
    mov x0, #1         // stdout
    mov x1, sp         // Use stack for newline
    mov w2, #10        // ASCII newline  
    strb w2, [x1]      // Store on stack
    mov x2, #1         // length
    svc #0x80
    
    add sp, sp, #32
    ldp x29, x30, [sp], #16
    ret

// Runtime function to print strings
.p2align 2
_print_string:
    stp x29, x30, [sp, #-16]!
    mov x29, sp
    
    // Calculate string length (simple strlen implementation)
    mov x1, x0         // Save string pointer
    mov x2, #0         // Length counter
    
.p2align 2
strlen_loop:
    ldrb w3, [x1, x2]  // Load byte at position
    cbz w3, print_str  // If null terminator, go to print
    add x2, x2, #1     // Increment length
    b strlen_loop
    
.p2align 2
print_str:
    // Write system call with calculated length
    mov x16, #4        // sys_write
    mov x0, #1         // stdout
    // x1 already has string pointer, x2 has length
    svc #0x80
    
    // Print newline
    mov x16, #4        // sys_write
    mov x0, #1         // stdout
    mov x1, sp         // Use stack for newline
    mov w2, #10        // ASCII newline
    strb w2, [x1]      // Store on stack
    mov x2, #1         // length
    svc #0x80
    
    ldp x29, x30, [sp], #16
    ret
`
	return runtime
}
