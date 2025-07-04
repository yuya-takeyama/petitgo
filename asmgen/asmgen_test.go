package asmgen

import (
	"strings"
	"testing"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/token"
)

func TestAsmGenerator_Generate_SimpleMain(t *testing.T) {
	gen := NewAsmGenerator()

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

	// Check basic assembly structure
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

func TestAsmGenerator_Generate_FibonacciFunction(t *testing.T) {
	gen := NewAsmGenerator()

	// Create fibonacci function: func fibonacci(n int) int { ... }
	param := ast.Parameter{
		Name: "n",
		Type: "int",
	}

	// if n <= 1 { return n }
	nVar := &ast.VariableNode{Name: "n"}
	one := &ast.NumberNode{Value: 1}
	condition := &ast.BinaryOpNode{
		Left:     nVar,
		Operator: token.LEQ,
		Right:    one,
	}
	returnStmt := &ast.ReturnStatement{Value: &ast.VariableNode{Name: "n"}}
	thenBlock := &ast.BlockStatement{Statements: []ast.Statement{returnStmt}}

	ifStmt := &ast.IfStatement{
		Condition: condition,
		ThenBlock: thenBlock,
	}

	// return fibonacci(n-1) + fibonacci(n-2)
	nMinus1 := &ast.BinaryOpNode{
		Left:     &ast.VariableNode{Name: "n"},
		Operator: token.SUB,
		Right:    &ast.NumberNode{Value: 1},
	}
	nMinus2 := &ast.BinaryOpNode{
		Left:     &ast.VariableNode{Name: "n"},
		Operator: token.SUB,
		Right:    &ast.NumberNode{Value: 2},
	}

	fib1 := &ast.CallNode{
		Function:  "fibonacci",
		Arguments: []ast.ASTNode{nMinus1},
	}
	fib2 := &ast.CallNode{
		Function:  "fibonacci",
		Arguments: []ast.ASTNode{nMinus2},
	}

	fibSum := &ast.BinaryOpNode{
		Left:     fib1,
		Operator: token.ADD,
		Right:    fib2,
	}

	finalReturn := &ast.ReturnStatement{Value: fibSum}

	funcBody := &ast.BlockStatement{
		Statements: []ast.Statement{ifStmt, finalReturn},
	}

	fibFunc := &ast.FuncStatement{
		Name:       "fibonacci",
		Parameters: []ast.Parameter{param},
		Body:       funcBody,
	}

	statements := []ast.Statement{fibFunc}
	result := gen.Generate(statements)

	// Check fibonacci function structure
	if !strings.Contains(result, "_fibonacci:") {
		t.Error("Missing fibonacci function label")
	}
	if !strings.Contains(result, "// Parameter: n") {
		t.Error("Missing parameter handling")
	}
	if !strings.Contains(result, "str x0, [x29, #-8]") {
		t.Error("Missing parameter storage")
	}
	if !strings.Contains(result, "bl _fibonacci") {
		t.Error("Missing recursive function call")
	}
	if !strings.Contains(result, "add x0, x1, x0") {
		t.Error("Missing addition operation")
	}
}

func TestAsmGenerator_Generate_VariableAssignment(t *testing.T) {
	gen := NewAsmGenerator()

	// Create: x := 42
	assignStmt := &ast.AssignStatement{
		Name:  "x",
		Value: &ast.NumberNode{Value: 42},
	}

	// Create main function with assignment
	blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt}}
	funcStmt := &ast.FuncStatement{
		Name:       "main",
		Parameters: []ast.Parameter{},
		Body:       blockStmt,
	}

	statements := []ast.Statement{funcStmt}
	result := gen.Generate(statements)

	// Check variable assignment
	if !strings.Contains(result, "mov x0, #42") {
		t.Error("Missing value loading")
	}
	if !strings.Contains(result, "str x0, [x29, #-8]") {
		t.Error("Missing variable storage")
	}
	if !strings.Contains(result, "// x := value") {
		t.Error("Missing assignment comment")
	}
}

func TestAsmGenerator_Generate_BinaryOperations(t *testing.T) {
	gen := NewAsmGenerator()

	tests := []struct {
		name     string
		op       token.Token
		expected string
	}{
		{"addition", token.ADD, "add x0, x1, x0"},
		{"subtraction", token.SUB, "sub x0, x1, x0"},
		{"multiplication", token.MUL, "mul x0, x1, x0"},
		{"division", token.QUO, "udiv x0, x1, x0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create binary operation: 5 op 3
			binOp := &ast.BinaryOpNode{
				Left:     &ast.NumberNode{Value: 5},
				Operator: tt.op,
				Right:    &ast.NumberNode{Value: 3},
			}

			printCall := &ast.CallNode{
				Function:  "println",
				Arguments: []ast.ASTNode{binOp},
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

			if !strings.Contains(result, tt.expected) {
				t.Errorf("Missing operation instruction: %s", tt.expected)
			}
		})
	}
}

func TestAsmGenerator_GenerateRuntime(t *testing.T) {
	gen := NewAsmGenerator()
	runtime := gen.GenerateRuntime()

	// Check that runtime functions are present
	expectedFunctions := []string{
		"_print_number:",
		"_print_char:",
		"print_nonzero:",
		"print_single_digit:",
		"print_done:",
	}

	for _, fn := range expectedFunctions {
		if !strings.Contains(runtime, fn) {
			t.Errorf("Missing runtime function: %s", fn)
		}
	}

	// Check specific implementation details
	if !strings.Contains(runtime, "mov x0, #48") {
		t.Error("Missing ASCII '0' conversion")
	}
	if !strings.Contains(runtime, "udiv x0, x0, x1") {
		t.Error("Missing division operation")
	}
	if !strings.Contains(runtime, "msub x0, x2, x1, x0") {
		t.Error("Missing modulo operation")
	}
	if !strings.Contains(runtime, "svc #0x80") {
		t.Error("Missing system call")
	}
}

func TestAsmGenerator_Generate_IfStatement(t *testing.T) {
	gen := NewAsmGenerator()

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
	exprStmt := &ast.ExpressionStatement{Expression: printCall}
	thenBlock := &ast.BlockStatement{Statements: []ast.Statement{exprStmt}}

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

	// Check if statement structure
	if !strings.Contains(result, "cmp x1, x0") {
		t.Error("Missing comparison instruction")
	}
	if !strings.Contains(result, "cset x0, hi") {
		t.Error("Missing condition evaluation")
	}
	if !strings.Contains(result, "beq L") {
		t.Error("Missing conditional branch")
	}
}

func TestAsmGenerator_Generate_ForLoop(t *testing.T) {
	gen := NewAsmGenerator()

	// Create: for i < 5 { println(i); i = i + 1 }
	condition := &ast.BinaryOpNode{
		Left:     &ast.VariableNode{Name: "i"},
		Operator: token.LSS,
		Right:    &ast.NumberNode{Value: 5},
	}

	printCall := &ast.CallNode{
		Function:  "println",
		Arguments: []ast.ASTNode{&ast.VariableNode{Name: "i"}},
	}
	exprStmt := &ast.ExpressionStatement{Expression: printCall}

	increment := &ast.BinaryOpNode{
		Left:     &ast.VariableNode{Name: "i"},
		Operator: token.ADD,
		Right:    &ast.NumberNode{Value: 1},
	}
	reassignStmt := &ast.ReassignStatement{
		Name:  "i",
		Value: increment,
	}

	forBody := &ast.BlockStatement{Statements: []ast.Statement{exprStmt, reassignStmt}}
	forStmt := &ast.ForStatement{
		Condition: condition,
		Body:      forBody,
	}

	initStmt := &ast.AssignStatement{
		Name:  "i",
		Value: &ast.NumberNode{Value: 0},
	}

	blockStmt := &ast.BlockStatement{Statements: []ast.Statement{initStmt, forStmt}}
	funcStmt := &ast.FuncStatement{
		Name:       "main",
		Parameters: []ast.Parameter{},
		Body:       blockStmt,
	}

	statements := []ast.Statement{funcStmt}
	result := gen.Generate(statements)

	// Check for loop structure
	if !strings.Contains(result, "beq L") {
		t.Error("Missing loop exit branch")
	}
	if !strings.Contains(result, "b L") {
		t.Error("Missing loop back branch")
	}
}
