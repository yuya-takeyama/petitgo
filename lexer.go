package main

type TokenType string

const (
	NUMBER TokenType = "NUMBER"
	PLUS   TokenType = "PLUS"
	MINUS  TokenType = "MINUS"
	STAR   TokenType = "STAR"
	SLASH  TokenType = "SLASH"
	EOF    TokenType = "EOF"
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
	// 入力の終端チェック
	if l.position >= len(l.input) {
		return Token{Type: EOF, Literal: ""}
	}
	
	ch := l.input[l.position]
	
	// 演算子の判定
	switch ch {
	case '+':
		l.position++
		return Token{Type: PLUS, Literal: "+"}
	case '-':
		l.position++
		return Token{Type: MINUS, Literal: "-"}
	case '*':
		l.position++
		return Token{Type: STAR, Literal: "*"}
	case '/':
		l.position++
		return Token{Type: SLASH, Literal: "/"}
	}
	
	// 数字を読み取る
	if isDigit(ch) {
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
	
	// 未知の文字
	l.position++
	return Token{Type: EOF, Literal: ""}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}