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
	input    string
	position int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, position: 0}
}

func (l *Lexer) NextToken() Token {
	// 数字を読み取る
	start := l.position
	for l.position < len(l.input) && isDigit(l.input[l.position]) {
		l.position++
	}
	
	literal := l.input[start:l.position]
	
	return Token{
		Type:    NUMBER,
		Literal: literal,
	}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}