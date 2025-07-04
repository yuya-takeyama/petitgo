package scanner

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/token"
)

func TestLexer_NextToken_SingleNumber(t *testing.T) {
	input := "123"

	scanner := NewScanner(input)

	tok := scanner.NextToken()

	if tok.Type != token.INT {
		t.Fatalf("token type wrong. expected=%d, got=%d", token.INT, tok.Type)
	}

	if tok.Literal != "123" {
		t.Fatalf("token literal wrong. expected=%s, got=%s", "123", tok.Literal)
	}
}

func TestLexer_NextToken_DifferentNumbers(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
	}{
		{"456", "456"},
		{"1", "1"},
		{"99999", "99999"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != token.INT {
			t.Fatalf("token type wrong. expected=%d, got=%d", token.INT, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_Operators(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"+", token.ADD, "+"},
		{"-", token.SUB, "-"},
		{"*", token.MUL, "*"},
		{"/", token.QUO, "/"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_Parentheses(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"(", token.LPAREN, "("},
		{")", token.RPAREN, ")"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_WhitespaceSkipping(t *testing.T) {
	input := "  123  +  456  "

	expected := []struct {
		expectedType    token.Token
		expectedLiteral string
	}{
		{token.INT, "123"},
		{token.ADD, "+"},
		{token.INT, "456"},
		{token.EOF, ""},
	}

	scanner := NewScanner(input)

	for i, tt := range expected {
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - token type wrong. expected=%d, got=%d", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - token literal wrong. expected=%s, got=%s", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_ComplexExpression(t *testing.T) {
	input := "1 + 2 * (3 - 4) / 5"

	expected := []struct {
		expectedType    token.Token
		expectedLiteral string
	}{
		{token.INT, "1"},
		{token.ADD, "+"},
		{token.INT, "2"},
		{token.MUL, "*"},
		{token.LPAREN, "("},
		{token.INT, "3"},
		{token.SUB, "-"},
		{token.INT, "4"},
		{token.RPAREN, ")"},
		{token.QUO, "/"},
		{token.INT, "5"},
		{token.EOF, ""},
	}

	scanner := NewScanner(input)

	for i, tt := range expected {
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - token type wrong. expected=%d, got=%d", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - token literal wrong. expected=%s, got=%s", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_ComparisonOperators(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"==", token.EQL, "=="},
		{"!=", token.NEQ, "!="},
		{"<", token.LSS, "<"},
		{">", token.GTR, ">"},
		{"<=", token.LEQ, "<="},
		{">=", token.GEQ, ">="},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_LogicalOperators(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"&&", token.LAND, "&&"},
		{"||", token.LOR, "||"},
		{"!", token.NOT, "!"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_Comma(t *testing.T) {
	input := ","
	scanner := NewScanner(input)
	tok := scanner.NextToken()

	if tok.Type != token.COMMA {
		t.Fatalf("token type wrong. expected=%d, got=%d", token.COMMA, tok.Type)
	}

	if tok.Literal != "," {
		t.Fatalf("token literal wrong. expected=%s, got=%s", ",", tok.Literal)
	}
}

func TestLexer_NextToken_Braces(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"{", token.LBRACE, "{"},
		{"}", token.RBRACE, "}"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_Keywords(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"if", token.IF, "if"},
		{"else", token.ELSE, "else"},
		{"for", token.FOR, "for"},
		{"break", token.BREAK, "break"},
		{"continue", token.CONTINUE, "continue"},
		{"func", token.FUNC, "func"},
		{"return", token.RETURN, "return"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_ControlFlowExpression(t *testing.T) {
	input := "if x > 5 { print(x) } else { print(0) }"

	expected := []struct {
		expectedType    token.Token
		expectedLiteral string
	}{
		{token.IF, "if"},
		{token.IDENT, "x"},
		{token.GTR, ">"},
		{token.INT, "5"},
		{token.LBRACE, "{"},
		{token.IDENT, "print"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.IDENT, "print"},
		{token.LPAREN, "("},
		{token.INT, "0"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	scanner := NewScanner(input)

	for i, tt := range expected {
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - token type wrong. expected=%d, got=%d", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - token literal wrong. expected=%s, got=%s", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
