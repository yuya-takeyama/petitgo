package main

import "testing"

func TestEvalNumber(t *testing.T) {
	node := &NumberNode{Value: 42}
	result := Eval(node)

	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func TestEvalAddition(t *testing.T) {
	left := &NumberNode{Value: 1}
	right := &NumberNode{Value: 2}
	node := &BinaryOpNode{
		Left:     left,
		Operator: ADD,
		Right:    right,
	}

	result := Eval(node)

	if result != 3 {
		t.Errorf("expected 3, got %d", result)
	}
}

func TestEvalSubtraction(t *testing.T) {
	left := &NumberNode{Value: 10}
	right := &NumberNode{Value: 5}
	node := &BinaryOpNode{
		Left:     left,
		Operator: SUB,
		Right:    right,
	}

	result := Eval(node)

	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

func TestEvalMultiplication(t *testing.T) {
	left := &NumberNode{Value: 3}
	right := &NumberNode{Value: 4}
	node := &BinaryOpNode{
		Left:     left,
		Operator: MUL,
		Right:    right,
	}

	result := Eval(node)

	if result != 12 {
		t.Errorf("expected 12, got %d", result)
	}
}

func TestEvalDivision(t *testing.T) {
	left := &NumberNode{Value: 8}
	right := &NumberNode{Value: 2}
	node := &BinaryOpNode{
		Left:     left,
		Operator: QUO,
		Right:    right,
	}

	result := Eval(node)

	if result != 4 {
		t.Errorf("expected 4, got %d", result)
	}
}

func TestEvalComplexExpression(t *testing.T) {
	// (2 + 3) * 4 = 20
	addNode := &BinaryOpNode{
		Left:     &NumberNode{Value: 2},
		Operator: ADD,
		Right:    &NumberNode{Value: 3},
	}

	mulNode := &BinaryOpNode{
		Left:     addNode,
		Operator: MUL,
		Right:    &NumberNode{Value: 4},
	}

	result := Eval(mulNode)

	if result != 20 {
		t.Errorf("expected 20, got %d", result)
	}
}

func TestEvalWithPrecedence(t *testing.T) {
	// 2 + 3 * 4 = 14 (not 20)
	mulNode := &BinaryOpNode{
		Left:     &NumberNode{Value: 3},
		Operator: MUL,
		Right:    &NumberNode{Value: 4},
	}

	addNode := &BinaryOpNode{
		Left:     &NumberNode{Value: 2},
		Operator: ADD,
		Right:    mulNode,
	}

	result := Eval(addNode)

	if result != 14 {
		t.Errorf("expected 14, got %d", result)
	}
}

// 統合テスト: Lexer → Parser → Evaluator
func TestEvalIntegration(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"42", 42},
		{"1 + 2", 3},
		{"10 - 5", 5},
		{"3 * 4", 12},
		{"8 / 2", 4},
		{"2 + 3 * 4", 14},
		{"(2 + 3) * 4", 20},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		parser := NewParser(lexer)
		ast := parser.ParseExpression()
		result := Eval(ast)

		if result != tt.expected {
			t.Errorf("input: %s, expected: %d, got: %d", tt.input, tt.expected, result)
		}
	}
}
