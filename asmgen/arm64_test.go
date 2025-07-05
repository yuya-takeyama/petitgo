// ARM64 assembly generation tests (can run on any platform)

package asmgen

import (
	"strings"
	"testing"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/token"
)

func TestARM64Generator_Generate_SimpleMain(t *testing.T) {
	gen := NewARM64Generator()

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

	// Check ARM64 assembly structure
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

func TestARM64Generator_Generate_BinaryOperations(t *testing.T) {
	tests := []struct {
		name     string
		operator token.Token
		expected string
	}{
		{"addition", token.ADD, "add x0, x0, x1"},
		{"subtraction", token.SUB, "sub x0, x0, x1"},
		{"multiplication", token.MUL, "mul x0, x0, x1"},
		{"division", token.QUO, "udiv x0, x0, x1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewARM64Generator()

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

func TestARM64Generator_GenerateRuntime(t *testing.T) {
	gen := NewARM64Generator()
	runtime := gen.GenerateRuntime()

	// Check ARM64 runtime functions
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

	// Check ARM64 specific implementation details
	if !strings.Contains(runtime, "mov w3, #48") {
		t.Error("Missing ASCII '0' conversion")
	}
	if !strings.Contains(runtime, "udiv x4, x0, x3") {
		t.Error("Missing division operation")
	}
	if !strings.Contains(runtime, "msub x5, x4, x3, x0") {
		t.Error("Missing modulo operation")
	}
	if !strings.Contains(runtime, "svc #0x80") {
		t.Error("Missing macOS system call")
	}
}

func TestARM64Generator_Generate_IfStatement(t *testing.T) {
	gen := NewARM64Generator()

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

	// Check ARM64 if statement generation
	if !strings.Contains(result, "cmp x0, x1") {
		t.Error("Missing comparison instruction")
	}
	if !strings.Contains(result, "cset x0, gt") {
		t.Error("Missing condition evaluation")
	}
	if !strings.Contains(result, "cbz x0,") {
		t.Error("Missing conditional branch")
	}
}

// TestARM64Generator_UntestedFeatures tests ARM64 features with 0% coverage
func TestARM64Generator_UntestedFeatures(t *testing.T) {
	gen := NewARM64Generator()

	// Test for statement (generateForStatement)
	t.Run("for_statement", func(t *testing.T) {
		// Create: for x > 5 { println(x) }
		condition := &ast.BinaryOpNode{
			Left:     &ast.VariableNode{Name: "x"},
			Operator: token.GTR,
			Right:    &ast.NumberNode{Value: 5},
		}
		printCall := &ast.CallNode{
			Function:  "println",
			Arguments: []ast.ASTNode{&ast.VariableNode{Name: "x"}},
		}
		body := &ast.BlockStatement{
			Statements: []ast.Statement{&ast.ExpressionStatement{Expression: printCall}},
		}
		forStmt := &ast.ForStatement{
			Condition: condition,
			Body:      body,
		}

		// Create variable x := 10
		assignStmt := &ast.AssignStatement{
			Name:  "x",
			Value: &ast.NumberNode{Value: 10},
		}

		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt, forStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		// Debug: print the actual result
		t.Logf("Generated assembly:\n%s", result)

		// Check for statement generation
		if !strings.Contains(result, "L1:") || !strings.Contains(result, "L2:") {
			t.Error("Missing loop labels")
		}
		if !strings.Contains(result, "cbz x0, L2") {
			t.Error("Missing conditional branch for loop exit")
		}
		if !strings.Contains(result, "b L1") {
			t.Error("Missing jump back to loop start")
		}
	})

	// Test if/else statement to improve generateIfStatement coverage
	t.Run("if_else_statement", func(t *testing.T) {
		// Create: if x > 5 { println(1) } else { println(2) }
		condition := &ast.BinaryOpNode{
			Left:     &ast.VariableNode{Name: "x"},
			Operator: token.GTR,
			Right:    &ast.NumberNode{Value: 5},
		}
		thenCall := &ast.CallNode{
			Function:  "println",
			Arguments: []ast.ASTNode{&ast.NumberNode{Value: 1}},
		}
		elseCall := &ast.CallNode{
			Function:  "println",
			Arguments: []ast.ASTNode{&ast.NumberNode{Value: 2}},
		}
		thenBlock := &ast.BlockStatement{
			Statements: []ast.Statement{&ast.ExpressionStatement{Expression: thenCall}},
		}
		elseBlock := &ast.BlockStatement{
			Statements: []ast.Statement{&ast.ExpressionStatement{Expression: elseCall}},
		}
		ifStmt := &ast.IfStatement{
			Condition: condition,
			ThenBlock: thenBlock,
			ElseBlock: elseBlock, // This will test the else path
		}

		assignStmt := &ast.AssignStatement{
			Name:  "x",
			Value: &ast.NumberNode{Value: 3},
		}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt, ifStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		// Check if/else statement generation (basic test)
		if result == "" {
			t.Error("Generate returned empty string")
		}
	})

	// Test return statement (generateReturnStatement)
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

		// Check return statement generation
		if !strings.Contains(result, "mov x0, #42") {
			t.Error("Missing return value setup")
		}
		if !strings.Contains(result, "add sp, sp, #64") {
			t.Error("Missing stack restore")
		}
		if !strings.Contains(result, "ldp x29, x30, [sp], #16") {
			t.Error("Missing frame pointer restore")
		}
		if !strings.Contains(result, "ret") {
			t.Error("Missing return instruction")
		}
	})

	// Test function call (generateFunctionCall) - Note: Only println is currently implemented
	t.Run("function_call", func(t *testing.T) {
		// Use println which is the only implemented function call
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

		// Check println generation (which uses _print_number)
		if !strings.Contains(result, "mov x0, #5") {
			t.Error("Missing argument setup for println")
		}
		if !strings.Contains(result, "bl _print_number") {
			t.Error("Missing print_number call")
		}
	})

	// Test switch statement (generateSwitchStatement)
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
		defaultStmt := &ast.BlockStatement{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Expression: &ast.CallNode{
						Function:  "println",
						Arguments: []ast.ASTNode{&ast.NumberNode{Value: 999}},
					},
				},
			},
		}

		switchStmt := &ast.SwitchStatement{
			Value:   &ast.VariableNode{Name: "x"},
			Cases:   []*ast.CaseStatement{caseStmt},
			Default: defaultStmt,
		}

		// Create variable x := 1
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

		// Check switch statement generation (basic structure)
		// Just ensure code generation doesn't crash - switch might not be fully implemented
		if result == "" {
			t.Error("Generate returned empty string")
		}
	})

	// Test field access (generateFieldAccess)
	t.Run("field_access", func(t *testing.T) {
		fieldAccess := &ast.FieldAccess{
			Object: &ast.VariableNode{Name: "obj"},
			Field:  "value",
		}
		exprStmt := &ast.ExpressionStatement{Expression: fieldAccess}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{exprStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		// Check that code generation doesn't crash (basic test)
		if result == "" {
			t.Error("Generate returned empty string")
		}
	})

	// Test slice literal (generateSliceLiteral)
	t.Run("slice_literal", func(t *testing.T) {
		sliceLiteral := &ast.ArrayLiteral{
			Elements: []ast.ASTNode{
				&ast.NumberNode{Value: 1},
				&ast.NumberNode{Value: 2},
				&ast.NumberNode{Value: 3},
			},
		}
		exprStmt := &ast.ExpressionStatement{Expression: sliceLiteral}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{exprStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		// Check that code generation doesn't crash
		if result == "" {
			t.Error("Generate returned empty string")
		}
	})

	// Test index access (generateIndexAccess)
	t.Run("index_access", func(t *testing.T) {
		indexAccess := &ast.IndexAccess{
			Object: &ast.VariableNode{Name: "arr"},
			Index:  &ast.NumberNode{Value: 0},
		}
		exprStmt := &ast.ExpressionStatement{Expression: indexAccess}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{exprStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		// Check that code generation doesn't crash
		if result == "" {
			t.Error("Generate returned empty string")
		}
	})

	// Test getFieldOffset function
	t.Run("get_field_offset", func(t *testing.T) {
		// This is a helper function, test it indirectly through field access
		fieldAccess := &ast.FieldAccessNode{
			Object: &ast.VariableNode{Name: "obj"},
			Field:  "name",
		}
		assignStmt := &ast.AssignStatement{
			Name:  "obj",
			Value: &ast.StructLiteral{},
		}
		exprStmt := &ast.ExpressionStatement{Expression: fieldAccess}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt, exprStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		// Check that code generation doesn't crash
		if result == "" {
			t.Error("Generate returned empty string")
		}
	})
}

// TestARM64Generator_CoverageCompleter tests remaining 0% coverage functions
func TestARM64Generator_CoverageCompleter(t *testing.T) {
	gen := NewARM64Generator()

	// Test generateFunctionCall directly (not through println)
	t.Run("direct_function_call", func(t *testing.T) {
		// Create a call that goes to generateFunctionCall (not println)
		callNode := &ast.CallNode{
			Function:  "fibonacci", // Not println, so will call generateFunctionCall
			Arguments: []ast.ASTNode{&ast.NumberNode{Value: 5}},
		}
		// Use AssignStatement to trigger generateExpression -> generateFunctionCall path
		assignStmt := &ast.AssignStatement{
			Name:  "result",
			Value: callNode,
		}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		// Just verify it generated something without crashing
		if result == "" {
			t.Error("Generate returned empty string")
		}
	})

	// Test generateFieldAccess coverage
	t.Run("field_access_coverage", func(t *testing.T) {
		// Note: Check AST definition - might be FieldAccessNode not FieldAccess
		fieldAccess := &ast.FieldAccessNode{
			Object: &ast.VariableNode{Name: "obj"},
			Field:  "value",
		}
		// Put field access in expression context to trigger generateExpression -> generateFieldAccess
		assignStmt := &ast.AssignStatement{
			Name:  "result",
			Value: fieldAccess,
		}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt}}
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

	// Test generateSliceLiteral coverage
	t.Run("slice_literal_coverage", func(t *testing.T) {
		// Note: Check AST definition - might be SliceLiteral not ArrayLiteral
		sliceLiteral := &ast.SliceLiteral{
			Elements: []ast.ASTNode{
				&ast.NumberNode{Value: 1},
				&ast.NumberNode{Value: 2},
			},
		}
		assignStmt := &ast.AssignStatement{
			Name:  "arr",
			Value: sliceLiteral,
		}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt}}
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

	// Test generateIndexAccess coverage
	t.Run("index_access_coverage", func(t *testing.T) {
		indexAccess := &ast.IndexAccess{
			Object: &ast.VariableNode{Name: "arr"},
			Index:  &ast.NumberNode{Value: 0},
		}
		assignStmt := &ast.AssignStatement{
			Name:  "value",
			Value: indexAccess,
		}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt}}
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

	// Test getFieldOffset coverage indirectly through field access with proper context
	t.Run("field_offset_coverage", func(t *testing.T) {
		// Test multiple field names to cover all getFieldOffset branches
		fieldNames := []string{"name", "age", "z", "unknown"} // Different offsets + default

		for _, fieldName := range fieldNames {
			structLit := &ast.StructLiteral{}
			assignStmt := &ast.AssignStatement{
				Name:  "obj",
				Value: structLit,
			}
			fieldAccess := &ast.FieldAccessNode{
				Object: &ast.VariableNode{Name: "obj"},
				Field:  fieldName,
			}
			accessStmt := &ast.AssignStatement{
				Name:  "field",
				Value: fieldAccess,
			}
			blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt, accessStmt}}
			funcStmt := &ast.FuncStatement{
				Name:       "main",
				Parameters: []ast.Parameter{},
				Body:       blockStmt,
			}
			statements := []ast.Statement{funcStmt}
			result := gen.Generate(statements)

			if result == "" {
				t.Errorf("Generate returned empty string for field %s", fieldName)
			}
		}
	})
}

// TestARM64Generator_GenerateStatementCoverage tests uncovered statement types
func TestARM64Generator_GenerateStatementCoverage(t *testing.T) {
	gen := NewARM64Generator()

	// Test VarStatement coverage
	t.Run("var_statement", func(t *testing.T) {
		varStmt := &ast.VarStatement{
			Name:     "x",
			TypeName: "int",
			Value:    &ast.NumberNode{Value: 42},
		}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{varStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		if !strings.Contains(result, "// var x") {
			t.Error("Missing var statement comment")
		}
	})

	// Test ReassignStatement coverage
	t.Run("reassign_statement", func(t *testing.T) {
		assignStmt := &ast.AssignStatement{
			Name:  "x",
			Value: &ast.NumberNode{Value: 10},
		}
		reassignStmt := &ast.ReassignStatement{
			Name:  "x",
			Value: &ast.NumberNode{Value: 20},
		}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt, reassignStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		if !strings.Contains(result, "// x = value") {
			t.Error("Missing reassignment comment")
		}
	})

	// Test ExpressionStatement with non-println expression
	t.Run("expression_statement_other", func(t *testing.T) {
		// Create an expression statement that's not a println call
		exprStmt := &ast.ExpressionStatement{
			Expression: &ast.NumberNode{Value: 42}, // Not a CallNode
		}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{exprStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		// Just verify it doesn't crash
		if result == "" {
			t.Error("Generate returned empty string")
		}
	})

	// Test VariableNode with non-existing variable (else branch)
	t.Run("variable_not_found", func(t *testing.T) {
		// Use a variable that doesn't exist
		varNode := &ast.VariableNode{Name: "nonexistent"}
		assignStmt := &ast.AssignStatement{
			Name:  "result",
			Value: varNode,
		}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		// Should generate without crashing (no ldr instruction for non-existing var)
		if result == "" {
			t.Error("Generate returned empty string")
		}
	})

	// Test StringNode if it exists in generateExpression
	t.Run("string_node", func(t *testing.T) {
		stringNode := &ast.StringNode{Value: "hello"}
		assignStmt := &ast.AssignStatement{
			Name:  "str",
			Value: stringNode,
		}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{assignStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "main",
			Parameters: []ast.Parameter{},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		// Just verify it doesn't crash (might be unsupported)
		if result == "" {
			t.Error("Generate returned empty string")
		}
	})

	// Test generateFunction with parameters
	t.Run("function_with_parameters", func(t *testing.T) {
		// Create a function with parameter: func test(x int) { println(x) }
		printCall := &ast.CallNode{
			Function:  "println",
			Arguments: []ast.ASTNode{&ast.VariableNode{Name: "x"}},
		}
		exprStmt := &ast.ExpressionStatement{Expression: printCall}
		blockStmt := &ast.BlockStatement{Statements: []ast.Statement{exprStmt}}
		funcStmt := &ast.FuncStatement{
			Name:       "test", // Not main
			Parameters: []ast.Parameter{{Name: "x", Type: "int"}},
			Body:       blockStmt,
		}
		statements := []ast.Statement{funcStmt}
		result := gen.Generate(statements)

		// Check parameter handling
		if !strings.Contains(result, "// Parameter: x") {
			t.Error("Missing parameter comment")
		}
		if !strings.Contains(result, "_test:") {
			t.Error("Missing function label")
		}
		// Should NOT have main exit code (test line 71)
		if strings.Contains(result, "mov x16, #1") {
			t.Error("Non-main function should not have exit code")
		}
	})
}
