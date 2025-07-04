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