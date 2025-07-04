package asmgen

import (
	"fmt"
	"strings"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/token"
)

// X86_64Generator generates x86_64 assembly code for Linux
type X86_64Generator struct {
	output    strings.Builder
	labelNum  int
	variables map[string]int // variable name -> stack offset
	stackSize int
}

// NewX86_64Generator creates a new x86_64 assembly generator
func NewX86_64Generator() *X86_64Generator {
	return &X86_64Generator{
		variables: make(map[string]int),
		stackSize: 0,
	}
}

// Generate converts an AST to x86_64 assembly code
func (g *X86_64Generator) Generate(statements []ast.Statement) string {
	g.output.Reset()

	// Basic ELF sections for Linux x86_64
	g.writeLine(".section .text")
	g.writeLine(".globl _start")
	g.writeLine("")

	// Generate all functions first
	for _, stmt := range statements {
		if funcStmt, ok := stmt.(*ast.FuncStatement); ok {
			g.generateFunction(funcStmt)
		}
	}

	return g.output.String()
}

func (g *X86_64Generator) generateFunction(funcStmt *ast.FuncStatement) {
	// Reset variables for each function
	g.variables = make(map[string]int)
	g.stackSize = 0

	// Linux uses _start as entry point instead of main
	if funcStmt.Name == "main" {
		g.writeLine("_start:")
	} else {
		g.writeLine(fmt.Sprintf("_%s:", funcStmt.Name))
	}

	g.writeLine("    # Function prologue")
	g.writeLine("    pushq %rbp")      // Save base pointer
	g.writeLine("    movq %rsp, %rbp") // Set base pointer
	g.writeLine("    subq $64, %rsp")  // Reserve 64 bytes for local variables

	// For functions with parameters, handle the first parameter in %rdi (Linux calling convention)
	if len(funcStmt.Parameters) > 0 {
		// Store parameter on stack
		g.stackSize += 8
		g.variables[funcStmt.Parameters[0].Name] = g.stackSize
		g.writeLine(fmt.Sprintf("    # Parameter: %s", funcStmt.Parameters[0].Name))
		g.writeLine(fmt.Sprintf("    movq %%rdi, -%d(%%rbp)", g.stackSize))
	}

	// Generate function body
	g.generateBlock(funcStmt.Body)

	// Function epilogue only for main (other functions use explicit return)
	if funcStmt.Name == "main" {
		g.writeLine("    # Exit")
		g.writeLine("    addq $64, %rsp") // Restore stack space
		g.writeLine("    popq %rbp")      // Restore base pointer
		g.writeLine("    movq $60, %rax") // sys_exit
		g.writeLine("    movq $0, %rdi")  // exit status
		g.writeLine("    syscall")        // system call
	}
	g.writeLine("")
}

func (g *X86_64Generator) generateBlock(block *ast.BlockStatement) {
	for _, stmt := range block.Statements {
		g.generateStatement(stmt)
	}
}

func (g *X86_64Generator) generateStatement(stmt ast.Statement) {
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
			g.writeLine(fmt.Sprintf("    # %s = value (reassignment)", s.Name))
			g.generateExpression(s.Value)
			g.writeLine(fmt.Sprintf("    movq %%rax, -%d(%%rbp)", offset))
		} else {
			// New variable assignment
			g.stackSize += 8
			g.variables[s.Name] = g.stackSize
			g.writeLine(fmt.Sprintf("    # %s := value", s.Name))
			g.generateExpression(s.Value)
			g.writeLine(fmt.Sprintf("    movq %%rax, -%d(%%rbp)", g.stackSize))
		}
	case *ast.VarStatement:
		g.stackSize += 8
		g.variables[s.Name] = g.stackSize
		g.writeLine(fmt.Sprintf("    # var %s", s.Name))
		g.generateExpression(s.Value)
		g.writeLine(fmt.Sprintf("    movq %%rax, -%d(%%rbp)", g.stackSize))
	case *ast.IfStatement:
		g.generateIfStatement(s)
	case *ast.ForStatement:
		g.generateForStatement(s)
	case *ast.ReturnStatement:
		g.generateReturnStatement(s)
	case *ast.ReassignStatement:
		// Variable reassignment
		if offset, exists := g.variables[s.Name]; exists {
			g.writeLine(fmt.Sprintf("    # %s = value", s.Name))
			g.generateExpression(s.Value)
			g.writeLine(fmt.Sprintf("    movq %%rax, -%d(%%rbp)", offset))
		}
	case *ast.SwitchStatement:
		g.generateSwitchStatement(s)
	}
}

func (g *X86_64Generator) generatePrintln(arg ast.ASTNode) {
	g.generateExpression(arg)
	g.writeLine("    # Print number in %rax")
	g.writeLine("    call _print_number")
}

func (g *X86_64Generator) generateExpression(expr ast.ASTNode) {
	switch e := expr.(type) {
	case *ast.NumberNode:
		g.writeLine(fmt.Sprintf("    movq $%d, %%rax", e.Value))
	case *ast.VariableNode:
		if offset, exists := g.variables[e.Name]; exists {
			g.writeLine(fmt.Sprintf("    movq -%d(%%rbp), %%rax", offset))
		}
	case *ast.BinaryOpNode:
		// Left operand
		g.generateExpression(e.Left)
		g.writeLine("    pushq %rax")

		// Right operand
		g.generateExpression(e.Right)
		g.writeLine("    movq %rax, %rbx")
		g.writeLine("    popq %rax")

		// Operation
		switch e.Operator {
		case token.ADD:
			g.writeLine("    addq %rbx, %rax")
		case token.SUB:
			g.writeLine("    subq %rbx, %rax")
		case token.MUL:
			g.writeLine("    imulq %rbx, %rax")
		case token.QUO:
			g.writeLine("    cqo")        // Sign extend %rax to %rdx:%rax
			g.writeLine("    idivq %rbx") // Divide %rdx:%rax by %rbx
		case token.EQL:
			g.writeLine("    cmpq %rbx, %rax")
			g.writeLine("    sete %al")
			g.writeLine("    movzbq %al, %rax")
		case token.NEQ:
			g.writeLine("    cmpq %rbx, %rax")
			g.writeLine("    setne %al")
			g.writeLine("    movzbq %al, %rax")
		case token.LSS:
			g.writeLine("    cmpq %rbx, %rax")
			g.writeLine("    setl %al")
			g.writeLine("    movzbq %al, %rax")
		case token.LEQ:
			g.writeLine("    cmpq %rbx, %rax")
			g.writeLine("    setle %al")
			g.writeLine("    movzbq %al, %rax")
		case token.GTR:
			g.writeLine("    cmpq %rbx, %rax")
			g.writeLine("    setg %al")
			g.writeLine("    movzbq %al, %rax")
		case token.GEQ:
			g.writeLine("    cmpq %rbx, %rax")
			g.writeLine("    setge %al")
			g.writeLine("    movzbq %al, %rax")
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

func (g *X86_64Generator) generateIfStatement(stmt *ast.IfStatement) {
	endLabel := g.getNewLabel()

	g.generateExpression(stmt.Condition)
	g.writeLine("    testq %rax, %rax")
	g.writeLine("    jz " + endLabel)

	g.generateBlock(stmt.ThenBlock)

	if stmt.ElseBlock != nil {
		elseLabel := g.getNewLabel()
		g.writeLine("    jmp " + elseLabel)
		g.writeLine(endLabel + ":")
		g.generateBlock(stmt.ElseBlock)
		g.writeLine(elseLabel + ":")
	} else {
		g.writeLine(endLabel + ":")
	}
}

func (g *X86_64Generator) generateForStatement(stmt *ast.ForStatement) {
	startLabel := g.getNewLabel()
	endLabel := g.getNewLabel()

	g.writeLine(startLabel + ":")
	g.generateExpression(stmt.Condition)
	g.writeLine("    testq %rax, %rax")
	g.writeLine("    jz " + endLabel)

	g.generateBlock(stmt.Body)

	g.writeLine("    jmp " + startLabel)
	g.writeLine(endLabel + ":")
}

func (g *X86_64Generator) generateReturnStatement(stmt *ast.ReturnStatement) {
	if stmt.Value != nil {
		g.generateExpression(stmt.Value)
	}
	g.writeLine("    addq $64, %rsp") // Restore stack space
	g.writeLine("    popq %rbp")      // Restore base pointer
	g.writeLine("    ret")            // Return
}

func (g *X86_64Generator) generateSwitchStatement(stmt *ast.SwitchStatement) {
	endLabel := g.getNewLabel()

	// Generate the switch value and store it
	g.writeLine("    # switch expression")
	g.generateExpression(stmt.Value)
	g.writeLine("    pushq %rax") // Store switch value on stack

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
		g.writeLine(fmt.Sprintf("    # case %d comparison", i))
		g.generateExpression(caseStmt.Value)
		g.writeLine("    movq (%rsp), %rbx") // Load switch value from stack
		g.writeLine("    cmpq %rax, %rbx")
		g.writeLine(fmt.Sprintf("    je %s", caseLabels[i]))
	}

	// If no case matched, go to default or end
	if stmt.Default != nil {
		g.writeLine(fmt.Sprintf("    jmp %s", defaultLabel))
	} else {
		g.writeLine(fmt.Sprintf("    jmp %s", endLabel))
	}

	// Generate case bodies
	for i, caseStmt := range stmt.Cases {
		g.writeLine(fmt.Sprintf("%s:", caseLabels[i]))
		g.writeLine(fmt.Sprintf("    # case %d body", i))
		g.generateBlock(caseStmt.Body)
		g.writeLine(fmt.Sprintf("    jmp %s", endLabel)) // break to end
	}

	// Generate default case if exists
	if stmt.Default != nil {
		g.writeLine(fmt.Sprintf("%s:", defaultLabel))
		g.writeLine("    # default case")
		g.generateBlock(stmt.Default)
	}

	// End label
	g.writeLine(fmt.Sprintf("%s:", endLabel))
	g.writeLine("    addq $8, %rsp") // Restore stack (remove stored switch value)
}

func (g *X86_64Generator) generateFunctionCall(call *ast.CallNode) {
	if len(call.Arguments) > 0 {
		g.generateExpression(call.Arguments[0])
		g.writeLine("    movq %rax, %rdi") // First argument goes to %rdi
	}
	g.writeLine(fmt.Sprintf("    call _%s", call.Function))
}

func (g *X86_64Generator) generateFieldAccess(node *ast.FieldAccessNode) {
	g.writeLine("    # Field access: obj.field")
	g.generateExpression(node.Object)
	fieldOffset := getFieldOffset(node.Field)
	g.writeLine(fmt.Sprintf("    # Access field '%s' at offset %d", node.Field, fieldOffset))
	g.writeLine(fmt.Sprintf("    movq %d(%%rax), %%rax", fieldOffset))
}

func (g *X86_64Generator) generateSliceLiteral(node *ast.SliceLiteral) {
	g.writeLine("    # Slice literal creation")
	elementCount := len(node.Elements)
	g.writeLine(fmt.Sprintf("    movq $%d, %%rax", elementCount))
}

func (g *X86_64Generator) generateIndexAccess(node *ast.IndexAccess) {
	g.writeLine("    # Index access: arr[index]")
	g.generateExpression(node.Object)
	g.writeLine("    pushq %rax") // Store array base address
	g.generateExpression(node.Index)
	g.writeLine("    popq %rbx")         // Load array base address
	g.writeLine("    salq $3, %rax")     // %rax = index * 8
	g.writeLine("    addq %rbx, %rax")   // %rax = base + offset
	g.writeLine("    movq (%rax), %rax") // Load value at address
}

func (g *X86_64Generator) getNewLabel() string {
	g.labelNum++
	return fmt.Sprintf("L%d", g.labelNum)
}

func (g *X86_64Generator) writeLine(s string) {
	g.output.WriteString(s + "\n")
}

func (g *X86_64Generator) GenerateRuntime() string {
	runtime := `
# Runtime function to print numbers (x86_64 Linux)
_print_number:
    pushq %rbp
    movq %rsp, %rbp
    subq $32, %rsp
    
    # Save the number
    movq %rax, -8(%rbp)
    
    # Buffer for digits (16 bytes should be enough)
    leaq -24(%rbp), %rsi  # Buffer pointer
    movq $0, %rcx         # Digit count
    
    # Handle zero case
    cmpq $0, %rax
    jne convert_loop
    movb $48, (%rsi)      # ASCII '0'
    movq $1, %rcx
    jmp print_digits
    
convert_loop:
    testq %rax, %rax
    jz print_digits
    
    # Divide by 10
    movq $10, %rbx
    cqo                   # Sign extend %rax to %rdx:%rax
    idivq %rbx            # %rax = quotient, %rdx = remainder
    
    # Convert digit to ASCII and store
    addq $48, %rdx        # Convert to ASCII
    movb %dl, (%rsi,%rcx,1) # Store digit
    incq %rcx             # Increment digit count
    
    jmp convert_loop
    
print_digits:
    # Print digits in reverse order
    testq %rcx, %rcx
    jz print_newline
    
print_loop:
    decq %rcx
    movb (%rsi,%rcx,1), %dl
    
    # Write system call
    movq $1, %rax         # sys_write
    movq $1, %rdi         # stdout
    leaq -25(%rbp), %rsi  # temporary buffer
    movb %dl, (%rsi)      # store character
    movq $1, %rdx         # length
    syscall
    
    testq %rcx, %rcx
    jnz print_loop
    
print_newline:
    # Write newline character
    movq $1, %rax         # sys_write
    movq $1, %rdi         # stdout
    leaq -25(%rbp), %rsi  # buffer
    movb $10, (%rsi)      # ASCII newline  
    movq $1, %rdx         # length
    syscall
    
    addq $32, %rsp
    popq %rbp
    ret
`
	return runtime
}
