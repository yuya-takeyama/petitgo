package token

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
