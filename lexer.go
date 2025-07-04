package main

// Token is the set of lexical tokens of the Go programming language.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	literal_beg
	// Identifiers and basic type literals
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	IMAG   // 123.45i
	CHAR   // 'a'
	STRING // "abc"
	literal_end

	operator_beg
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	ASSIGN // := or =

	// Comparison operators
	EQL // ==
	NEQ // !=
	LSS // <
	GTR // >
	LEQ // <=
	GEQ // >=

	// Logical operators
	LAND // &&
	LOR  // ||
	NOT  // !

	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }
	operator_end

	keyword_beg
	// Keywords
	IF       // if
	ELSE     // else
	FOR      // for
	BREAK    // break
	CONTINUE // continue
	keyword_end
)

type TokenInfo struct {
	Type    Token
	Literal string
}

// Keywords map for keyword detection
var keywords = map[string]Token{
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
	"var":      IDENT, // var is handled as IDENT for now
}

type Lexer struct {
	input    string
	position int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, position: 0}
}

func (l *Lexer) NextToken() TokenInfo {
	// 空白文字をスキップ
	l.skipWhitespace()

	// 入力の終端チェック
	if l.position >= len(l.input) {
		return TokenInfo{Type: EOF, Literal: ""}
	}

	ch := l.input[l.position]

	// 演算子と括弧の判定
	switch ch {
	case '+':
		l.position++
		return TokenInfo{Type: ADD, Literal: "+"}
	case '-':
		l.position++
		return TokenInfo{Type: SUB, Literal: "-"}
	case '*':
		l.position++
		return TokenInfo{Type: MUL, Literal: "*"}
	case '/':
		l.position++
		return TokenInfo{Type: QUO, Literal: "/"}
	case ':':
		if l.position+1 < len(l.input) && l.input[l.position+1] == '=' {
			l.position += 2
			return TokenInfo{Type: ASSIGN, Literal: ":="}
		}
	case '=':
		if l.position+1 < len(l.input) && l.input[l.position+1] == '=' {
			l.position += 2
			return TokenInfo{Type: EQL, Literal: "=="}
		}
		l.position++
		return TokenInfo{Type: ASSIGN, Literal: "="}
	case '!':
		if l.position+1 < len(l.input) && l.input[l.position+1] == '=' {
			l.position += 2
			return TokenInfo{Type: NEQ, Literal: "!="}
		}
		l.position++
		return TokenInfo{Type: NOT, Literal: "!"}
	case '<':
		if l.position+1 < len(l.input) && l.input[l.position+1] == '=' {
			l.position += 2
			return TokenInfo{Type: LEQ, Literal: "<="}
		}
		l.position++
		return TokenInfo{Type: LSS, Literal: "<"}
	case '>':
		if l.position+1 < len(l.input) && l.input[l.position+1] == '=' {
			l.position += 2
			return TokenInfo{Type: GEQ, Literal: ">="}
		}
		l.position++
		return TokenInfo{Type: GTR, Literal: ">"}
	case '&':
		if l.position+1 < len(l.input) && l.input[l.position+1] == '&' {
			l.position += 2
			return TokenInfo{Type: LAND, Literal: "&&"}
		}
	case '|':
		if l.position+1 < len(l.input) && l.input[l.position+1] == '|' {
			l.position += 2
			return TokenInfo{Type: LOR, Literal: "||"}
		}
	case '(':
		l.position++
		return TokenInfo{Type: LPAREN, Literal: "("}
	case ')':
		l.position++
		return TokenInfo{Type: RPAREN, Literal: ")"}
	case '{':
		l.position++
		return TokenInfo{Type: LBRACE, Literal: "{"}
	case '}':
		l.position++
		return TokenInfo{Type: RBRACE, Literal: "}"}
	}

	// 識別子を読み取る（文字で始まる）
	if isLetter(ch) {
		start := l.position
		for l.position < len(l.input) && (isLetter(l.input[l.position]) || isDigit(l.input[l.position])) {
			l.position++
		}

		literal := l.input[start:l.position]

		// キーワードかチェック
		if tokenType, isKeyword := keywords[literal]; isKeyword {
			return TokenInfo{
				Type:    tokenType,
				Literal: literal,
			}
		}

		return TokenInfo{
			Type:    IDENT,
			Literal: literal,
		}
	}

	// 数字を読み取る
	if isDigit(ch) {
		start := l.position
		for l.position < len(l.input) && isDigit(l.input[l.position]) {
			l.position++
		}

		literal := l.input[start:l.position]
		return TokenInfo{
			Type:    INT,
			Literal: literal,
		}
	}

	// 未知の文字
	l.position++
	return TokenInfo{Type: EOF, Literal: ""}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func (l *Lexer) skipWhitespace() {
	for l.position < len(l.input) && isWhitespace(l.input[l.position]) {
		l.position++
	}
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
