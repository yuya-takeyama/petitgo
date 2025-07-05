package scanner

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/token"
)

func TestLexer_NextToken_SingleNumber(t *testing.T) {
	input := "123"

	scanner := NewScanner(input)

	tok := scanner.NextToken()

	if tok.Type != token.INT {
		t.Fatalf("token type wrong. expected=%d, got=%d", token.INT, tok.Type)
	}

	if tok.Literal != "123" {
		t.Fatalf("token literal wrong. expected=%s, got=%s", "123", tok.Literal)
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
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != token.INT {
			t.Fatalf("token type wrong. expected=%d, got=%d", token.INT, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_Operators(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"+", token.ADD, "+"},
		{"-", token.SUB, "-"},
		{"*", token.MUL, "*"},
		{"/", token.QUO, "/"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_Parentheses(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"(", token.LPAREN, "("},
		{")", token.RPAREN, ")"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_WhitespaceSkipping(t *testing.T) {
	input := "  123  +  456  "

	expected := []struct {
		expectedType    token.Token
		expectedLiteral string
	}{
		{token.INT, "123"},
		{token.ADD, "+"},
		{token.INT, "456"},
		{token.EOF, ""},
	}

	scanner := NewScanner(input)

	for i, tt := range expected {
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - token type wrong. expected=%d, got=%d", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - token literal wrong. expected=%s, got=%s", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_ComplexExpression(t *testing.T) {
	input := "1 + 2 * (3 - 4) / 5"

	expected := []struct {
		expectedType    token.Token
		expectedLiteral string
	}{
		{token.INT, "1"},
		{token.ADD, "+"},
		{token.INT, "2"},
		{token.MUL, "*"},
		{token.LPAREN, "("},
		{token.INT, "3"},
		{token.SUB, "-"},
		{token.INT, "4"},
		{token.RPAREN, ")"},
		{token.QUO, "/"},
		{token.INT, "5"},
		{token.EOF, ""},
	}

	scanner := NewScanner(input)

	for i, tt := range expected {
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - token type wrong. expected=%d, got=%d", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - token literal wrong. expected=%s, got=%s", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_ComparisonOperators(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"==", token.EQL, "=="},
		{"!=", token.NEQ, "!="},
		{"<", token.LSS, "<"},
		{">", token.GTR, ">"},
		{"<=", token.LEQ, "<="},
		{">=", token.GEQ, ">="},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_LogicalOperators(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"&&", token.LAND, "&&"},
		{"||", token.LOR, "||"},
		{"!", token.NOT, "!"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_Comma(t *testing.T) {
	input := ","
	scanner := NewScanner(input)
	tok := scanner.NextToken()

	if tok.Type != token.COMMA {
		t.Fatalf("token type wrong. expected=%d, got=%d", token.COMMA, tok.Type)
	}

	if tok.Literal != "," {
		t.Fatalf("token literal wrong. expected=%s, got=%s", ",", tok.Literal)
	}
}

func TestLexer_NextToken_Braces(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"{", token.LBRACE, "{"},
		{"}", token.RBRACE, "}"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_Keywords(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"if", token.IF, "if"},
		{"else", token.ELSE, "else"},
		{"for", token.FOR, "for"},
		{"break", token.BREAK, "break"},
		{"continue", token.CONTINUE, "continue"},
		{"func", token.FUNC, "func"},
		{"return", token.RETURN, "return"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_NextToken_ControlFlowExpression(t *testing.T) {
	input := "if x > 5 { print(x) } else { print(0) }"

	expected := []struct {
		expectedType    token.Token
		expectedLiteral string
	}{
		{token.IF, "if"},
		{token.IDENT, "x"},
		{token.GTR, ">"},
		{token.INT, "5"},
		{token.LBRACE, "{"},
		{token.IDENT, "print"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.IDENT, "print"},
		{token.LPAREN, "("},
		{token.INT, "0"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	scanner := NewScanner(input)

	for i, tt := range expected {
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - token type wrong. expected=%d, got=%d", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - token literal wrong. expected=%s, got=%s", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test Input() method
func TestScanner_Input(t *testing.T) {
	input := "test input string"
	scanner := NewScanner(input)

	if scanner.Input() != input {
		t.Fatalf("Input() wrong. expected=%s, got=%s", input, scanner.Input())
	}
}

// Test Position() method
func TestScanner_Position(t *testing.T) {
	input := "123 + 456"
	scanner := NewScanner(input)

	// Initial position should be 0
	if scanner.Position() != 0 {
		t.Fatalf("Initial position wrong. expected=0, got=%d", scanner.Position())
	}

	// After scanning first token, position should advance
	scanner.NextToken() // "123"
	pos1 := scanner.Position()
	if pos1 <= 0 {
		t.Fatalf("Position should advance after scanning. got=%d", pos1)
	}

	scanner.NextToken() // "+"
	pos2 := scanner.Position()
	if pos2 <= pos1 {
		t.Fatalf("Position should advance further. pos1=%d, pos2=%d", pos1, pos2)
	}
}

// Test string literals
func TestScanner_StringLiterals(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{`"hello"`, token.STRING, "hello"},
		{`"world"`, token.STRING, "world"},
		{`""`, token.STRING, ""},
		{`"hello world"`, token.STRING, "hello world"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test string escape sequences
func TestScanner_StringEscapeSequences(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
	}{
		{`"hello\nworld"`, "hello\nworld"},
		{`"hello\tworld"`, "hello\tworld"},
		{`"hello\rworld"`, "hello\rworld"},
		{`"say \"hello\""`, `say "hello"`},
		{`"path\\file"`, `path\file`},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != token.STRING {
			t.Fatalf("token type wrong. expected=%d, got=%d", token.STRING, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test unclosed string literal
func TestScanner_UncloseStringLiteral(t *testing.T) {
	input := `"unclosed string`
	scanner := NewScanner(input)
	tok := scanner.NextToken()

	if tok.Type != token.ILLEGAL {
		t.Fatalf("token type wrong. expected=%d, got=%d", token.ILLEGAL, tok.Type)
	}

	if tok.Literal != "unclosed string" {
		t.Fatalf("token literal wrong. expected=%s, got=%s", "unclosed string", tok.Literal)
	}
}

// Test line comments
func TestScanner_LineComments(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
	}{
		{"// this is a comment", "// this is a comment"},
		{"// another comment", "// another comment"},
		{"//no space", "//no space"},
		{"// comment with spaces  ", "// comment with spaces  "},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != token.COMMENT {
			t.Fatalf("token type wrong. expected=%d, got=%d", token.COMMENT, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test block comments
func TestScanner_BlockComments(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
	}{
		{"/* single line */", "/* single line */"},
		{"/* multi\nline\ncomment */", "/* multi\nline\ncomment */"},
		{"/*no spaces*/", "/*no spaces*/"},
		{"/* nested // comment */", "/* nested // comment */"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != token.COMMENT {
			t.Fatalf("token type wrong. expected=%d, got=%d", token.COMMENT, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test increment and decrement operators
func TestScanner_IncrementDecrement(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"++", token.INC, "++"},
		{"--", token.DEC, "--"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test compound assignment operators
func TestScanner_CompoundAssignment(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"+=", token.ADD_ASSIGN, "+="},
		{"-=", token.SUB_ASSIGN, "-="},
		{"*=", token.MUL_ASSIGN, "*="},
		{"/=", token.QUO_ASSIGN, "/="},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test assignment operators
func TestScanner_AssignmentOperators(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{":=", token.ASSIGN, ":="},
		{"=", token.ASSIGN, "="},
		{":", token.COLON, ":"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test unknown escape sequences
func TestScanner_UnknownEscapeSequences(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
	}{
		{`"hello\xworld"`, `hello\xworld`}, // Unknown escape \x
		{`"test\zmore"`, `test\zmore`},     // Unknown escape \z
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != token.STRING {
			t.Fatalf("token type wrong. expected=%d, got=%d", token.STRING, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test additional punctuation tokens
func TestScanner_AdditionalPunctuation(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{".", token.PERIOD, "."},
		{";", token.SEMICOLON, ";"},
		{"[", token.LBRACK, "["},
		{"]", token.RBRACK, "]"},
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test single & and | characters (fallthrough cases)
func TestScanner_SingleAmpersandPipe(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Token
		expectedLiteral string
	}{
		{"&", token.EOF, ""}, // & alone falls through to EOF
		{"|", token.EOF, ""}, // | alone falls through to EOF
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("token type wrong. expected=%d, got=%d", tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. expected=%s, got=%s", tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test unknown characters
func TestScanner_UnknownCharacters(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"@"}, // Unknown character
		{"#"}, // Unknown character
		{"$"}, // Unknown character
		{"%"}, // Unknown character
		{"^"}, // Unknown character
		{"~"}, // Unknown character
		{"`"}, // Unknown character
	}

	for _, tt := range tests {
		scanner := NewScanner(tt.input)
		tok := scanner.NextToken()

		if tok.Type != token.EOF {
			t.Fatalf("token type wrong for unknown char %s. expected=%d, got=%d", tt.input, token.EOF, tok.Type)
		}

		if tok.Literal != "" {
			t.Fatalf("token literal wrong for unknown char %s. expected=\"\", got=%s", tt.input, tok.Literal)
		}
	}
}
