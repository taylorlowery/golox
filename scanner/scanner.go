package scanner

import (
	"strconv"

	"github.com/taylorlowery/lox/golox"
	"github.com/taylorlowery/lox/token"
)

var keywords = map[string]token.TokenType{
	"and":    token.AND,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"fun":    token.FUN,
	"for":    token.FOR,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}

type Scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		tokens: []token.Token{},
		line:   1,
	}
}

func (s Scanner) Source() string {
	return s.source
}

func (s Scanner) Tokens() []token.Token {
	return s.tokens
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) ScanTokens() []token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token.Token{TokenType: token.EOF, Lexeme: "", Literal: nil, Line: s.line})
	return s.tokens
}

func (s *Scanner) addToken(t token.TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.Token{TokenType: t, Lexeme: text, Literal: literal, Line: s.line})
}

func (s *Scanner) advance() byte {
	r := s.source[s.current]
	s.current++
	return r
}

// match is a conditional advance
// depending on whether the next character is a given expected character
func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.Source()[s.current] != byte(expected) {
		return false
	}
	s.current++
	return true
}

// string consumes the scanner source from one double quote to another
// TODO: support escape characters
func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		golox.Error(s.line, "unterminated string")
		return
	}

	s.advance()

	// gather string value without the quotes
	value := s.Source()[s.start+1 : s.current-1]
	s.addToken(token.STRING, value)
}

func (s Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.Source()[s.current]
}

func (s Scanner) peekNext() byte {
	if s.current+1 > len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// handle decimal
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
	}

	for isDigit(s.peek()) {
		s.advance()
	}

	numString := s.Source()[s.start:s.current]
	value, err := strconv.ParseFloat(numString, 64)
	if err != nil {
		golox.Error(s.line, "invalid number format")
		return
	}

	s.addToken(token.NUMBER, value)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	value := s.Source()[s.start:s.current]
	tokenType, ok := keywords[value]
	if !ok {
		tokenType = token.IDENTIFIER
	}
	s.addToken(tokenType, nil)
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func isAlpha(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || b == '_'
}

func isAlphaNumeric(b byte) bool {
	return isAlpha(b) || isDigit(b)
}

func (s *Scanner) scanToken() {
	var c byte = s.advance()
	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN, nil)
	case ')':
		s.addToken(token.RIGHT_PAREN, nil)
	case '{':
		s.addToken(token.LEFT_BRACE, nil)
	case '}':
		s.addToken(token.RIGHT_BRACE, nil)
	case ',':
		s.addToken(token.COMMA, nil)
	case '.':
		s.addToken(token.DOT, nil)
	case '-':
		s.addToken(token.MINUS, nil)
	case '+':
		s.addToken(token.PLUS, nil)
	case ';':
		s.addToken(token.SEMICOLON, nil)
	case '*':
		s.addToken(token.STAR, nil)
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL, nil)
		} else {
			s.addToken(token.BANG, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL, nil)
		} else {
			s.addToken(token.EQUAL, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL, nil)
		} else {
			s.addToken(token.LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL, nil)
		} else {
			s.addToken(token.GREATER, nil)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH, nil)
		}
	case ' ', '\r', '\t':
		// ignore whitespace
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			golox.Error(s.line, "unexpected character")
		}
	}
}
