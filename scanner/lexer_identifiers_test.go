package scanner

import "testing"

func TestLexerIdentifiers(t *testing.T) {
	input := "var x foo bar123"

	tests := []struct {
		expectedType    Token
		expectedLiteral string
	}{
		{IDENT, "var"},
		{IDENT, "x"},
		{IDENT, "foo"},
		{IDENT, "bar123"},
		{EOF, ""},
	}

	lexer := NewLexer(input)

	for i, tt := range tests {
		tok := lexer.NextToken()

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
		expectedType    Token
		expectedLiteral string
	}{
		{ASSIGN, ":="},
		{ASSIGN, "="},
		{EOF, ""},
	}

	lexer := NewLexer(input)

	for i, tt := range tests {
		tok := lexer.NextToken()

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
