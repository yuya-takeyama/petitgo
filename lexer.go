package main

type TokenType string

const (
	NUMBER TokenType = "NUMBER"
)

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input string
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) NextToken() Token {
	return Token{
		Type:    NUMBER,
		Literal: "123",
	}
}