package main

import "testing"

func TestParseNumber(t *testing.T) {
	input := "42"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	ast := parser.ParseExpression()

	numberNode, ok := ast.(*NumberNode)
	if !ok {
		t.Fatalf("expected NumberNode, got %T", ast)
	}

	if numberNode.Value != 42 {
		t.Errorf("expected 42, got %d", numberNode.Value)
	}
}

func TestParseBinaryExpression(t *testing.T) {
	input := "1 + 2"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	ast := parser.ParseExpression()

	binaryNode, ok := ast.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode, got %T", ast)
	}

	if binaryNode.Operator != ADD {
		t.Errorf("expected ADD operator, got %v", binaryNode.Operator)
	}

	leftNumber, ok := binaryNode.Left.(*NumberNode)
	if !ok {
		t.Fatalf("expected left operand to be NumberNode, got %T", binaryNode.Left)
	}
	if leftNumber.Value != 1 {
		t.Errorf("expected left operand to be 1, got %d", leftNumber.Value)
	}

	rightNumber, ok := binaryNode.Right.(*NumberNode)
	if !ok {
		t.Fatalf("expected right operand to be NumberNode, got %T", binaryNode.Right)
	}
	if rightNumber.Value != 2 {
		t.Errorf("expected right operand to be 2, got %d", rightNumber.Value)
	}
}

func TestParseOperatorPrecedence(t *testing.T) {
	input := "2 + 3 * 4"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	ast := parser.ParseExpression()

	// (2 + (3 * 4)) の構造になってるはず
	binaryNode, ok := ast.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode, got %T", ast)
	}

	if binaryNode.Operator != ADD {
		t.Errorf("expected top-level operator to be ADD, got %v", binaryNode.Operator)
	}

	// 右側が 3 * 4 の BinaryOpNode になってるはず
	rightBinary, ok := binaryNode.Right.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected right operand to be BinaryOpNode, got %T", binaryNode.Right)
	}

	if rightBinary.Operator != MUL {
		t.Errorf("expected right operator to be MUL, got %v", rightBinary.Operator)
	}
}

func TestParseParentheses(t *testing.T) {
	input := "(2 + 3) * 4"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	ast := parser.ParseExpression()

	// ((2 + 3) * 4) の構造になってるはず
	binaryNode, ok := ast.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode, got %T", ast)
	}

	if binaryNode.Operator != MUL {
		t.Errorf("expected top-level operator to be MUL, got %v", binaryNode.Operator)
	}

	// 左側が (2 + 3) の BinaryOpNode になってるはず
	leftBinary, ok := binaryNode.Left.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected left operand to be BinaryOpNode, got %T", binaryNode.Left)
	}

	if leftBinary.Operator != ADD {
		t.Errorf("expected left operator to be ADD, got %v", leftBinary.Operator)
	}
}

func TestParseComparison(t *testing.T) {
	tests := []struct {
		input    string
		operator Token
	}{
		{"x > 5", GTR},
		{"x < 5", LSS},
		{"x == 5", EQL},
		{"x != 5", NEQ},
		{"x >= 5", GEQ},
		{"x <= 5", LEQ},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		parser := NewParser(lexer)

		ast := parser.ParseExpression()

		binaryNode, ok := ast.(*BinaryOpNode)
		if !ok {
			t.Fatalf("expected BinaryOpNode, got %T", ast)
		}

		if binaryNode.Operator != tt.operator {
			t.Errorf("expected operator %v, got %v", tt.operator, binaryNode.Operator)
		}
	}
}

func TestParseIfStatement(t *testing.T) {
	input := "if x > 5 { print(x) }"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	stmt := parser.ParseStatement()

	ifStmt, ok := stmt.(*IfStatement)
	if !ok {
		t.Fatalf("expected IfStatement, got %T", stmt)
	}

	// condition should be x > 5
	binaryNode, ok := ifStmt.Condition.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected condition to be BinaryOpNode, got %T", ifStmt.Condition)
	}

	if binaryNode.Operator != GTR {
		t.Errorf("expected condition operator to be GTR, got %v", binaryNode.Operator)
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
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	stmt := parser.ParseStatement()

	ifStmt, ok := stmt.(*IfStatement)
	if !ok {
		t.Fatalf("expected IfStatement, got %T", stmt)
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
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	stmt := parser.ParseStatement()

	forStmt, ok := stmt.(*ForStatement)
	if !ok {
		t.Fatalf("expected ForStatement, got %T", stmt)
	}

	// condition should be x < 5
	binaryNode, ok := forStmt.Condition.(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected condition to be BinaryOpNode, got %T", forStmt.Condition)
	}

	if binaryNode.Operator != LSS {
		t.Errorf("expected condition operator to be LSS, got %v", binaryNode.Operator)
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
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	stmt := parser.ParseStatement()

	_, ok := stmt.(*BreakStatement)
	if !ok {
		t.Fatalf("expected BreakStatement, got %T", stmt)
	}
}

func TestParseContinueStatement(t *testing.T) {
	input := "continue"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	stmt := parser.ParseStatement()

	_, ok := stmt.(*ContinueStatement)
	if !ok {
		t.Fatalf("expected ContinueStatement, got %T", stmt)
	}
}

func TestParseBlockStatement(t *testing.T) {
	input := "{ x := 1 y := 2 }"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	stmt := parser.ParseStatement()

	blockStmt, ok := stmt.(*BlockStatement)
	if !ok {
		t.Fatalf("expected BlockStatement, got %T", stmt)
	}

	if len(blockStmt.Statements) != 2 {
		t.Errorf("expected 2 statements in block, got %d", len(blockStmt.Statements))
	}
}
