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
