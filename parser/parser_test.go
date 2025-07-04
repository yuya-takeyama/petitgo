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
