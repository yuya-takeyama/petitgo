package parser

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/scanner"
	"github.com/yuya-takeyama/petitgo/token"
)

func TestParseNumber(t *testing.T) {
	input := "42"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	expr := parser.ParseExpression()

	numberNode, ok := expr.(*ast.NumberNode)
	if !ok {
		t.Fatalf("expected *ast.NumberNode, got %T", expr)
	}

	if numberNode.Value != 42 {
		t.Errorf("expected 42, got %d", numberNode.Value)
	}
}

func TestParseBinaryExpression(t *testing.T) {
	input := "1 + 2"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	expr := parser.ParseExpression()

	binaryNode, ok := expr.(*ast.BinaryOpNode)
	if !ok {
		t.Fatalf("expected *ast.BinaryOpNode, got %T", expr)
	}

	if binaryNode.Operator != token.ADD {
		t.Errorf("expected token.ADD operator, got %v", binaryNode.Operator)
	}

	leftNumber, ok := binaryNode.Left.(*ast.NumberNode)
	if !ok {
		t.Fatalf("expected left operand to be ast.NumberNode, got %T", binaryNode.Left)
	}
	if leftNumber.Value != 1 {
		t.Errorf("expected left operand to be 1, got %d", leftNumber.Value)
	}

	rightNumber, ok := binaryNode.Right.(*ast.NumberNode)
	if !ok {
		t.Fatalf("expected right operand to be ast.NumberNode, got %T", binaryNode.Right)
	}
	if rightNumber.Value != 2 {
		t.Errorf("expected right operand to be 2, got %d", rightNumber.Value)
	}
}

func TestParseOperatorPrecedence(t *testing.T) {
	input := "2 + 3 * 4"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	expr := parser.ParseExpression()

	// (2 + (3 * 4)) の構造になってるはず
	binaryNode, ok := expr.(*ast.BinaryOpNode)
	if !ok {
		t.Fatalf("expected *ast.BinaryOpNode, got %T", expr)
	}

	if binaryNode.Operator != token.ADD {
		t.Errorf("expected top-level operator to be token.ADD, got %v", binaryNode.Operator)
	}

	// 右側が 3 * 4 の ast.BinaryOpNode になってるはず
	rightBinary, ok := binaryNode.Right.(*ast.BinaryOpNode)
	if !ok {
		t.Fatalf("expected right operand to be ast.BinaryOpNode, got %T", binaryNode.Right)
	}

	if rightBinary.Operator != token.MUL {
		t.Errorf("expected right operator to be token.MUL, got %v", rightBinary.Operator)
	}
}

func TestParseParentheses(t *testing.T) {
	input := "(2 + 3) * 4"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	expr := parser.ParseExpression()

	// ((2 + 3) * 4) の構造になってるはず
	binaryNode, ok := expr.(*ast.BinaryOpNode)
	if !ok {
		t.Fatalf("expected *ast.BinaryOpNode, got %T", expr)
	}

	if binaryNode.Operator != token.MUL {
		t.Errorf("expected top-level operator to be token.MUL, got %v", binaryNode.Operator)
	}

	// 左側が (2 + 3) の ast.BinaryOpNode になってるはず
	leftBinary, ok := binaryNode.Left.(*ast.BinaryOpNode)
	if !ok {
		t.Fatalf("expected left operand to be ast.BinaryOpNode, got %T", binaryNode.Left)
	}

	if leftBinary.Operator != token.ADD {
		t.Errorf("expected left operator to be token.ADD, got %v", leftBinary.Operator)
	}
}

func TestParseComparison(t *testing.T) {
	tests := []struct {
		input    string
		operator token.Token
	}{
		{"x > 5", token.GTR},
		{"x < 5", token.LSS},
		{"x == 5", token.EQL},
		{"x != 5", token.NEQ},
		{"x >= 5", token.GEQ},
		{"x <= 5", token.LEQ},
	}

	for _, tt := range tests {
		sc := scanner.NewScanner(tt.input)
		parser := NewParser(sc)

		expr := parser.ParseExpression()

		binaryNode, ok := expr.(*ast.BinaryOpNode)
		if !ok {
			t.Fatalf("expected *ast.BinaryOpNode, got %T", expr)
		}

		if binaryNode.Operator != tt.operator {
			t.Errorf("expected operator %v, got %v", tt.operator, binaryNode.Operator)
		}
	}
}

func TestParseIfStatement(t *testing.T) {
	input := "if x > 5 { print(x) }"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	ifStmt, ok := stmt.(*ast.IfStatement)
	if !ok {
		t.Fatalf("expected *ast.IfStatement, got %T", stmt)
	}

	// condition should be x > 5
	binaryNode, ok := ifStmt.Condition.(*ast.BinaryOpNode)
	if !ok {
		t.Fatalf("expected condition to be ast.BinaryOpNode, got %T", ifStmt.Condition)
	}

	if binaryNode.Operator != token.GTR {
		t.Errorf("expected condition operator to be token.GTR, got %v", binaryNode.Operator)
	}

	// then block should exist
	if ifStmt.ThenBlock == nil {
		t.Fatalf("expected ThenBlock to exist")
	}

	// else block should be nil
	if ifStmt.ElseBlock != nil {
		t.Errorf("expected ElseBlock to be nil")
	}
}

func TestParseIfElseStatement(t *testing.T) {
	input := "if x > 5 { print(x) } else { print(0) }"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	ifStmt, ok := stmt.(*ast.IfStatement)
	if !ok {
		t.Fatalf("expected *ast.IfStatement, got %T", stmt)
	}

	// then block should exist
	if ifStmt.ThenBlock == nil {
		t.Fatalf("expected ThenBlock to exist")
	}

	// else block should exist
	if ifStmt.ElseBlock == nil {
		t.Fatalf("expected ElseBlock to exist")
	}
}

func TestParseForStatement(t *testing.T) {
	input := "for x < 5 { print(x) }"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	forStmt, ok := stmt.(*ast.ForStatement)
	if !ok {
		t.Fatalf("expected *ast.ForStatement, got %T", stmt)
	}

	// condition should be x < 5
	binaryNode, ok := forStmt.Condition.(*ast.BinaryOpNode)
	if !ok {
		t.Fatalf("expected condition to be ast.BinaryOpNode, got %T", forStmt.Condition)
	}

	if binaryNode.Operator != token.LSS {
		t.Errorf("expected condition operator to be token.LSS, got %v", binaryNode.Operator)
	}

	// body should exist
	if forStmt.Body == nil {
		t.Fatalf("expected Body to exist")
	}

	// init and update should be nil for condition-only form
	if forStmt.Init != nil {
		t.Errorf("expected Init to be nil")
	}
	if forStmt.Update != nil {
		t.Errorf("expected Update to be nil")
	}
}

func TestParseBreakStatement(t *testing.T) {
	input := "break"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	_, ok := stmt.(*ast.BreakStatement)
	if !ok {
		t.Fatalf("expected *ast.BreakStatement, got %T", stmt)
	}
}

func TestParseContinueStatement(t *testing.T) {
	input := "continue"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	_, ok := stmt.(*ast.ContinueStatement)
	if !ok {
		t.Fatalf("expected *ast.ContinueStatement, got %T", stmt)
	}
}

func TestParseBlockStatement(t *testing.T) {
	input := "{ x := 1 y := 2 }"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	blockStmt, ok := stmt.(*ast.BlockStatement)
	if !ok {
		t.Fatalf("expected *ast.BlockStatement, got %T", stmt)
	}

	if len(blockStmt.Statements) != 2 {
		t.Errorf("expected 2 statements in block, got %d", len(blockStmt.Statements))
	}
}

// Test variable declaration (var statement)
func TestParseVarStatement(t *testing.T) {
	input := "var x int = 42"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	varStmt, ok := stmt.(*ast.VarStatement)
	if !ok {
		t.Fatalf("expected *ast.VarStatement, got %T", stmt)
	}

	if varStmt.Name != "x" {
		t.Errorf("expected variable name 'x', got %s", varStmt.Name)
	}

	if varStmt.TypeName != "int" {
		t.Errorf("expected variable type 'int', got %s", varStmt.TypeName)
	}

	numberNode, ok := varStmt.Value.(*ast.NumberNode)
	if !ok {
		t.Fatalf("expected value to be NumberNode, got %T", varStmt.Value)
	}

	if numberNode.Value != 42 {
		t.Errorf("expected value 42, got %d", numberNode.Value)
	}
}

// Test reassignment statement
func TestParseReassignStatement(t *testing.T) {
	input := "x = 10"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	reassignStmt, ok := stmt.(*ast.ReassignStatement)
	if !ok {
		t.Fatalf("expected *ast.ReassignStatement, got %T", stmt)
	}

	if reassignStmt.Name != "x" {
		t.Errorf("expected variable name 'x', got %s", reassignStmt.Name)
	}

	numberNode, ok := reassignStmt.Value.(*ast.NumberNode)
	if !ok {
		t.Fatalf("expected value to be NumberNode, got %T", reassignStmt.Value)
	}

	if numberNode.Value != 10 {
		t.Errorf("expected value 10, got %d", numberNode.Value)
	}
}

// Test increment statement
func TestParseIncStatement(t *testing.T) {
	input := "x++"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	incStmt, ok := stmt.(*ast.IncStatement)
	if !ok {
		t.Fatalf("expected *ast.IncStatement, got %T", stmt)
	}

	if incStmt.Name != "x" {
		t.Errorf("expected variable name 'x', got %s", incStmt.Name)
	}
}

// Test decrement statement
func TestParseDecStatement(t *testing.T) {
	input := "x--"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	decStmt, ok := stmt.(*ast.DecStatement)
	if !ok {
		t.Fatalf("expected *ast.DecStatement, got %T", stmt)
	}

	if decStmt.Name != "x" {
		t.Errorf("expected variable name 'x', got %s", decStmt.Name)
	}
}

// Test compound assignment statements
func TestParseCompoundAssignStatement(t *testing.T) {
	tests := []struct {
		input            string
		expectedOperator token.Token
	}{
		{"x += 5", token.ADD_ASSIGN},
		{"y -= 3", token.SUB_ASSIGN},
		{"z *= 2", token.MUL_ASSIGN},
		{"w /= 4", token.QUO_ASSIGN},
	}

	for _, tt := range tests {
		sc := scanner.NewScanner(tt.input)
		parser := NewParser(sc)

		stmt := parser.ParseStatement()

		compoundStmt, ok := stmt.(*ast.CompoundAssignStatement)
		if !ok {
			t.Fatalf("expected *ast.CompoundAssignStatement for %s, got %T", tt.input, stmt)
		}

		if compoundStmt.Operator != tt.expectedOperator {
			t.Errorf("expected operator %v for %s, got %v", tt.expectedOperator, tt.input, compoundStmt.Operator)
		}
	}
}

// Test switch statement
func TestParseSwitchStatement(t *testing.T) {
	input := `switch x {
	case 1:
		print("one")
	case 2:
		print("two")  
	default:
		print("other")
	}`
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	switchStmt, ok := stmt.(*ast.SwitchStatement)
	if !ok {
		t.Fatalf("expected *ast.SwitchStatement, got %T", stmt)
	}

	// Check value expression
	varNode, ok := switchStmt.Value.(*ast.VariableNode)
	if !ok {
		t.Fatalf("expected switch value to be VariableNode, got %T", switchStmt.Value)
	}
	if varNode.Name != "x" {
		t.Errorf("expected switch value 'x', got %s", varNode.Name)
	}

	// Check cases
	if len(switchStmt.Cases) != 2 {
		t.Errorf("expected 2 cases, got %d", len(switchStmt.Cases))
	}

	// Check default case exists
	if switchStmt.Default == nil {
		t.Errorf("expected default case to exist")
	}
}

// Test package statement
func TestParsePackageStatement(t *testing.T) {
	input := "package main"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	packageStmt, ok := stmt.(*ast.PackageStatement)
	if !ok {
		t.Fatalf("expected *ast.PackageStatement, got %T", stmt)
	}

	if packageStmt.Name != "main" {
		t.Errorf("expected package name 'main', got %s", packageStmt.Name)
	}
}

// Test import statement
func TestParseImportStatement(t *testing.T) {
	input := `import "fmt"`
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	importStmt, ok := stmt.(*ast.ImportStatement)
	if !ok {
		t.Fatalf("expected *ast.ImportStatement, got %T", stmt)
	}

	if importStmt.Path != "fmt" {
		t.Errorf("expected import path 'fmt', got %s", importStmt.Path)
	}
}

// Test full for statement with init, condition, and update
func TestParseFullForStatement(t *testing.T) {
	input := "for i := 0; i < 10; i++ { print(i) }"
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	forStmt, ok := stmt.(*ast.ForStatement)
	if !ok {
		t.Fatalf("expected *ast.ForStatement, got %T", stmt)
	}

	// Check init statement
	if forStmt.Init == nil {
		t.Fatalf("expected init statement to exist")
	}

	initStmt, ok := forStmt.Init.(*ast.AssignStatement)
	if !ok {
		t.Fatalf("expected init to be AssignStatement, got %T", forStmt.Init)
	}

	if initStmt.Name != "i" {
		t.Errorf("expected init variable 'i', got %s", initStmt.Name)
	}

	// Check condition
	if forStmt.Condition == nil {
		t.Fatalf("expected condition to exist")
	}

	// Check update statement
	if forStmt.Update == nil {
		t.Fatalf("expected update statement to exist")
	}

	updateStmt, ok := forStmt.Update.(*ast.IncStatement)
	if !ok {
		t.Fatalf("expected update to be IncStatement, got %T", forStmt.Update)
	}

	if updateStmt.Name != "i" {
		t.Errorf("expected update variable 'i', got %s", updateStmt.Name)
	}
}

// Test struct definition
func TestParseStructDefinition(t *testing.T) {
	input := `type Person struct {
		name string
		age int
	}`
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	stmt := parser.ParseStatement()

	typeStmt, ok := stmt.(*ast.TypeStatement)
	if !ok {
		t.Fatalf("expected *ast.TypeStatement, got %T", stmt)
	}

	if typeStmt.Name != "Person" {
		t.Errorf("expected struct name 'Person', got %s", typeStmt.Name)
	}

	if len(typeStmt.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(typeStmt.Fields))
	}

	// Check first field
	if typeStmt.Fields[0].Name != "name" {
		t.Errorf("expected first field name 'name', got %s", typeStmt.Fields[0].Name)
	}
	if typeStmt.Fields[0].Type != "string" {
		t.Errorf("expected first field type 'string', got %s", typeStmt.Fields[0].Type)
	}

	// Check second field
	if typeStmt.Fields[1].Name != "age" {
		t.Errorf("expected second field name 'age', got %s", typeStmt.Fields[1].Name)
	}
	if typeStmt.Fields[1].Type != "int" {
		t.Errorf("expected second field type 'int', got %s", typeStmt.Fields[1].Type)
	}
}

// Test struct literal
func TestParseStructLiteral(t *testing.T) {
	input := `Person{name: "John", age: 30}`
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	expr := parser.ParseExpression()

	structLit, ok := expr.(*ast.StructLiteral)
	if !ok {
		t.Fatalf("expected *ast.StructLiteral, got %T", expr)
	}

	if structLit.TypeName != "Person" {
		t.Errorf("expected struct type 'Person', got %s", structLit.TypeName)
	}

	if len(structLit.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(structLit.Fields))
	}

	// Check that name and age fields exist
	if _, exists := structLit.Fields["name"]; !exists {
		t.Errorf("expected 'name' field to exist")
	}

	if _, exists := structLit.Fields["age"]; !exists {
		t.Errorf("expected 'age' field to exist")
	}
}

// Test slice literal
func TestParseSliceLiteral(t *testing.T) {
	input := `[]int{1, 2, 3}`
	sc := scanner.NewScanner(input)
	parser := NewParser(sc)

	expr := parser.ParseExpression()

	sliceLit, ok := expr.(*ast.SliceLiteral)
	if !ok {
		t.Fatalf("expected *ast.SliceLiteral, got %T", expr)
	}

	if sliceLit.ElementType != "int" {
		t.Errorf("expected slice type 'int', got %s", sliceLit.ElementType)
	}

	if len(sliceLit.Elements) != 3 {
		t.Errorf("expected 3 elements, got %d", len(sliceLit.Elements))
	}

	// Check first element
	firstNum, ok := sliceLit.Elements[0].(*ast.NumberNode)
	if !ok {
		t.Fatalf("expected first element to be NumberNode, got %T", sliceLit.Elements[0])
	}
	if firstNum.Value != 1 {
		t.Errorf("expected first element value 1, got %d", firstNum.Value)
	}
}
