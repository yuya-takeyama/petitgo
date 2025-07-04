package main

import (
	"testing"
)

func TestLexer_NextToken_SingleNumber(t *testing.T) {
	input := "123"
	
	lexer := NewLexer(input)
	
	token := lexer.NextToken()
	
	if token.Type != NUMBER {
		t.Fatalf("token type wrong. expected=%s, got=%s", NUMBER, token.Type)
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
		
		if token.Type != NUMBER {
			t.Fatalf("token type wrong. expected=%s, got=%s", NUMBER, token.Type)
		}
		
		if token.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, token.Literal)
		}
	}
}

func TestLexer_NextToken_Operators(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    TokenType
		expectedLiteral string
	}{
		{"+", PLUS, "+"},
		{"-", MINUS, "-"},
		{"*", STAR, "*"},
		{"/", SLASH, "/"},
	}
	
	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		token := lexer.NextToken()
		
		if token.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%s, got=%s", tt.expectedType, token.Type)
		}
		
		if token.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, token.Literal)
		}
	}
}