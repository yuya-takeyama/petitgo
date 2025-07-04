package main

import (
	"testing"
)

func TestLexer_NextToken_SingleNumber(t *testing.T) {
	input := "123"

	lexer := NewLexer(input)

	token := lexer.NextToken()

	if token.Type != INT {
		t.Fatalf("token type wrong. expected=%d, got=%d", INT, token.Type)
	}

	if token.Literal != "123" {
		t.Fatalf("token literal wrong. expected=%s, got=%s", "123", token.Literal)
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
		lexer := NewLexer(tt.input)
		token := lexer.NextToken()

		if token.Type != INT {
			t.Fatalf("token type wrong. expected=%d, got=%d", INT, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, token.Literal)
		}
	}
}

func TestLexer_NextToken_Operators(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    Token
		expectedLiteral string
	}{
		{"+", ADD, "+"},
		{"-", SUB, "-"},
		{"*", MUL, "*"},
		{"/", QUO, "/"},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		token := lexer.NextToken()

		if token.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, token.Literal)
		}
	}
}

func TestLexer_NextToken_Parentheses(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    Token
		expectedLiteral string
	}{
		{"(", LPAREN, "("},
		{")", RPAREN, ")"},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		token := lexer.NextToken()

		if token.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, token.Literal)
		}
	}
}

func TestLexer_NextToken_WhitespaceSkipping(t *testing.T) {
	input := "  123  +  456  "

	expected := []struct {
		expectedType    Token
		expectedLiteral string
	}{
		{INT, "123"},
		{ADD, "+"},
		{INT, "456"},
		{EOF, ""},
	}

	lexer := NewLexer(input)

	for i, tt := range expected {
		token := lexer.NextToken()

		if token.Type != tt.expectedType {
			t.Fatalf("test[%d] - token type wrong. expected=%d, got=%d", i, tt.expectedType, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - token literal wrong. expected=%s, got=%s", i, tt.expectedLiteral, token.Literal)
		}
	}
}

func TestLexer_NextToken_ComplexExpression(t *testing.T) {
	input := "1 + 2 * (3 - 4) / 5"

	expected := []struct {
		expectedType    Token
		expectedLiteral string
	}{
		{INT, "1"},
		{ADD, "+"},
		{INT, "2"},
		{MUL, "*"},
		{LPAREN, "("},
		{INT, "3"},
		{SUB, "-"},
		{INT, "4"},
		{RPAREN, ")"},
		{QUO, "/"},
		{INT, "5"},
		{EOF, ""},
	}

	lexer := NewLexer(input)

	for i, tt := range expected {
		token := lexer.NextToken()

		if token.Type != tt.expectedType {
			t.Fatalf("test[%d] - token type wrong. expected=%d, got=%d", i, tt.expectedType, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - token literal wrong. expected=%s, got=%s", i, tt.expectedLiteral, token.Literal)
		}
	}
}
