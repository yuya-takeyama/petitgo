package scanner

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/token"
)

func TestLexerIdentifiers(t *testing.T) {
	input := "var x foo bar123"

	tests := []struct {
		expectedType    token.Token
		expectedLiteral string
	}{
		{token.IDENT, "var"},
		{token.IDENT, "x"},
		{token.IDENT, "foo"},
		{token.IDENT, "bar123"},
		{token.EOF, ""},
	}

	scanner := NewScanner(input)

	for i, tt := range tests {
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexerAssignmentOperators(t *testing.T) {
	input := ":= ="

	tests := []struct {
		expectedType    token.Token
		expectedLiteral string
	}{
		{token.ASSIGN, ":="},
		{token.ASSIGN, "="},
		{token.EOF, ""},
	}

	scanner := NewScanner(input)

	for i, tt := range tests {
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
