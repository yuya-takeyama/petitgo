package scanner

import "github.com/yuya-takeyama/petitgo/token"

// Keywords map for keyword detection
var keywords = map[string]token.Token{
	"if":       token.IF,
	"else":     token.ELSE,
	"for":      token.FOR,
	"break":    token.BREAK,
	"continue": token.CONTINUE,
	"func":     token.FUNC,
	"return":   token.RETURN,
	"true":     token.TRUE,
	"false":    token.FALSE,
	"var":      token.IDENT, // var is handled as token.IDENT for now
}

type Scanner struct {
	input    string
	position int
}

func NewScanner(input string) *Scanner {
	return &Scanner{input: input, position: 0}
}

func (s *Scanner) Input() string {
	return s.input
}

func (s *Scanner) Position() int {
	return s.position
}

func (s *Scanner) NextToken() token.TokenInfo {
	// Skip whitespace characters
	s.skipWhitespace()

	// Check for end of input
	if s.position >= len(s.input) {
		return token.TokenInfo{Type: token.EOF, Literal: ""}
	}

	ch := s.input[s.position]

	// Handle operators and delimiters
	switch ch {
	case '+':
		s.position++
		return token.TokenInfo{Type: token.ADD, Literal: "+"}
	case '-':
		s.position++
		return token.TokenInfo{Type: token.SUB, Literal: "-"}
	case '*':
		s.position++
		return token.TokenInfo{Type: token.MUL, Literal: "*"}
	case '/':
		s.position++
		return token.TokenInfo{Type: token.QUO, Literal: "/"}
	case ':':
		if s.position+1 < len(s.input) && s.input[s.position+1] == '=' {
			s.position += 2
			return token.TokenInfo{Type: token.ASSIGN, Literal: ":="}
		}
	case '=':
		if s.position+1 < len(s.input) && s.input[s.position+1] == '=' {
			s.position += 2
			return token.TokenInfo{Type: token.EQL, Literal: "=="}
		}
		s.position++
		return token.TokenInfo{Type: token.ASSIGN, Literal: "="}
	case '!':
		if s.position+1 < len(s.input) && s.input[s.position+1] == '=' {
			s.position += 2
			return token.TokenInfo{Type: token.NEQ, Literal: "!="}
		}
		s.position++
		return token.TokenInfo{Type: token.NOT, Literal: "!"}
	case '<':
		if s.position+1 < len(s.input) && s.input[s.position+1] == '=' {
			s.position += 2
			return token.TokenInfo{Type: token.LEQ, Literal: "<="}
		}
		s.position++
		return token.TokenInfo{Type: token.LSS, Literal: "<"}
	case '>':
		if s.position+1 < len(s.input) && s.input[s.position+1] == '=' {
			s.position += 2
			return token.TokenInfo{Type: token.GEQ, Literal: ">="}
		}
		s.position++
		return token.TokenInfo{Type: token.GTR, Literal: ">"}
	case '&':
		if s.position+1 < len(s.input) && s.input[s.position+1] == '&' {
			s.position += 2
			return token.TokenInfo{Type: token.LAND, Literal: "&&"}
		}
	case '|':
		if s.position+1 < len(s.input) && s.input[s.position+1] == '|' {
			s.position += 2
			return token.TokenInfo{Type: token.LOR, Literal: "||"}
		}
	case '(':
		s.position++
		return token.TokenInfo{Type: token.LPAREN, Literal: "("}
	case ')':
		s.position++
		return token.TokenInfo{Type: token.RPAREN, Literal: ")"}
	case '{':
		s.position++
		return token.TokenInfo{Type: token.LBRACE, Literal: "{"}
	case '}':
		s.position++
		return token.TokenInfo{Type: token.RBRACE, Literal: "}"}
	case ',':
		s.position++
		return token.TokenInfo{Type: token.COMMA, Literal: ","}
	case '"':
		return s.readString()
	}

	// Read identifier (starts with letter)
	if isLetter(ch) {
		start := s.position
		for s.position < len(s.input) && (isLetter(s.input[s.position]) || isDigit(s.input[s.position])) {
			s.position++
		}

		literal := s.input[start:s.position]

		// Check if it's a keyword
		if tokenType, isKeyword := keywords[literal]; isKeyword {
			return token.TokenInfo{
				Type:    tokenType,
				Literal: literal,
			}
		}

		return token.TokenInfo{
			Type:    token.IDENT,
			Literal: literal,
		}
	}

	// Read number
	if isDigit(ch) {
		start := s.position
		for s.position < len(s.input) && isDigit(s.input[s.position]) {
			s.position++
		}

		literal := s.input[start:s.position]
		return token.TokenInfo{
			Type:    token.INT,
			Literal: literal,
		}
	}

	// Unknown character
	s.position++
	return token.TokenInfo{Type: token.EOF, Literal: ""}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func (s *Scanner) skipWhitespace() {
	for s.position < len(s.input) && isWhitespace(s.input[s.position]) {
		s.position++
	}
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

// readString reads a string literal from input
func (s *Scanner) readString() token.TokenInfo {
	start := s.position
	s.position++ // skip opening quote

	for s.position < len(s.input) && s.input[s.position] != '"' {
		// Handle escape sequences
		if s.input[s.position] == '\\' && s.position+1 < len(s.input) {
			s.position += 2 // skip escape sequence
		} else {
			s.position++
		}
	}

	if s.position >= len(s.input) {
		// Unclosed string literal - return error token
		return token.TokenInfo{Type: token.ILLEGAL, Literal: "unclosed string"}
	}

	// Extract string content (without quotes)
	literal := s.input[start+1 : s.position]
	s.position++ // skip closing quote

	// Process escape sequences
	processed := processEscapeSequences(literal)

	return token.TokenInfo{Type: token.STRING, Literal: processed}
}

// processEscapeSequences handles basic escape sequences
func processEscapeSequences(input string) string {
	result := ""
	for i := 0; i < len(input); i++ {
		if input[i] == '\\' && i+1 < len(input) {
			switch input[i+1] {
			case 'n':
				result += "\n"
			case 't':
				result += "\t"
			case 'r':
				result += "\r"
			case '"':
				result += "\""
			case '\\':
				result += "\\"
			default:
				// Unknown escape sequence, keep as is
				result += string(input[i])
				result += string(input[i+1])
			}
			i++ // skip next character
		} else {
			result += string(input[i])
		}
	}
	return result
}
