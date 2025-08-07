package scanner_test

import (
	"strings"
	"testing"

	"github.com/taylorlowery/lox/scanner"
	"github.com/taylorlowery/lox/token"
)

func TestScanner_SingleCharacterTokens_CorrectlyScanned(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		source    string
		output    string
		expectErr bool
	}{
		{
			name:      "parentheses",
			source:    "()",
			output:    "LEFT_PAREN RIGHT_PAREN EOF",
			expectErr: false,
		},
		{
			name:      "braces",
			source:    "{}",
			output:    "LEFT_BRACE RIGHT_BRACE EOF",
			expectErr: false,
		},
		{
			name:      "comma",
			source:    ",",
			output:    "COMMA EOF",
			expectErr: false,
		},
		{
			name:      "dot",
			source:    ".",
			output:    "DOT EOF",
			expectErr: false,
		},
		{
			name:      "minus",
			source:    "-",
			output:    "MINUS EOF",
			expectErr: false,
		},
		{
			name:      "plus",
			source:    "+",
			output:    "PLUS EOF",
			expectErr: false,
		},
		{
			name:      "semicolon",
			source:    ";",
			output:    "SEMICOLON EOF",
			expectErr: false,
		},
		{
			name:      "star",
			source:    "*",
			output:    "STAR EOF",
			expectErr: false,
		},
		{
			name:      "all single characters",
			source:    "(){},.-+;*",
			output:    "LEFT_PAREN RIGHT_PAREN LEFT_BRACE RIGHT_BRACE COMMA DOT MINUS PLUS SEMICOLON STAR EOF",
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := scanner.NewScanner(tc.source)
			if s.Source() != tc.source {
				t.Fatalf("expected new scanner to have source %q, got %q", tc.source, s.Source())
			}
			if len(s.Tokens()) > 0 {
				t.Fatalf("expected new scanner to have empty token list, got %#v", s)
			}
			got := s.ScanTokens()
			// Log tokens for debugging
			t.Logf("Source: %q", tc.source)
			for i, token := range got {
				t.Logf("[%d] %s", i, token.TokenType.String())
			}
			tokenTypes := tokenListAsString(got)
			if tokenTypes != tc.output {
				t.Fatalf("expected token types %q, got %q", tc.output, tokenTypes)
			}
		})
	}
}

func tokenListAsString(tokens []token.Token) string {
	tokenTypes := make([]string, len(tokens))
	for i, t := range tokens {
		tokenTypes[i] = t.TokenType.String()
	}
	return strings.Join(tokenTypes, " ")
}

func TestScanner_EqualSignTokens_CorrectlyScanned(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		source    string
		output    string
		expectErr bool
	}{
		{
			name:      "single equal",
			source:    "=",
			output:    "EQUAL EOF",
			expectErr: false,
		},
		{
			name:      "double equal",
			source:    "==",
			output:    "EQUAL_EQUAL EOF",
			expectErr: false,
		},
		{
			name:      "bang equal",
			source:    "!=",
			output:    "BANG_EQUAL EOF",
			expectErr: false,
		},
		{
			name:      "single bang",
			source:    "!",
			output:    "BANG EOF",
			expectErr: false,
		},
		{
			name:      "less than equal",
			source:    "<=",
			output:    "LESS_EQUAL EOF",
			expectErr: false,
		},
		{
			name:      "less than",
			source:    "<",
			output:    "LESS EOF",
			expectErr: false,
		},
		{
			name:      "greater than equal",
			source:    ">=",
			output:    "GREATER_EQUAL EOF",
			expectErr: false,
		},
		{
			name:      "greater than",
			source:    ">",
			output:    "GREATER EOF",
			expectErr: false,
		},
		{
			name:      "mixed comparison operators",
			source:    "!===<=>=",
			output:    "BANG_EQUAL EQUAL_EQUAL LESS_EQUAL GREATER_EQUAL EOF",
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := scanner.NewScanner(tc.source)
			if s.Source() != tc.source {
				t.Fatalf("expected new scanner to have source %q, got %q", tc.source, s.Source())
			}
			if len(s.Tokens()) > 0 {
				t.Fatalf("expected new scanner to have empty token list, got %#v", s)
			}
			got := s.ScanTokens()
			// Log tokens for debugging
			t.Logf("Source: %q", tc.source)
			for i, token := range got {
				t.Logf("[%d] %s", i, token.TokenType.String())
			}
			tokenTypes := tokenListAsString(got)
			if tokenTypes != tc.output {
				t.Fatalf("expected token types %q, got %q", tc.output, tokenTypes)
			}
		})
	}
}

func TestScanner_SlashTokens_CorrectlyScanned(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		source    string
		output    string
		expectErr bool
	}{
		{
			name:      "single slash",
			source:    "/",
			output:    "SLASH EOF",
			expectErr: false,
		},
		{
			name:      "slash in expression",
			source:    "10/5",
			output:    "SLASH EOF", // Note: numbers not implemented yet
			expectErr: true,        // Will error on digits
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := scanner.NewScanner(tc.source)
			got := s.ScanTokens()
			// Log tokens for debugging
			t.Logf("Source: %q", tc.source)
			for i, token := range got {
				t.Logf("[%d] %s", i, token.TokenType.String())
			}
			if !tc.expectErr {
				tokenTypes := tokenListAsString(got)
				if tokenTypes != tc.output {
					t.Fatalf("expected token types %q, got %q", tc.output, tokenTypes)
				}
			}
		})
	}
}

func TestScanner_Whitespace_IgnoredCorrectly(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		source string
		output string
	}{
		{
			name:   "spaces between tokens",
			source: "( )",
			output: "LEFT_PAREN RIGHT_PAREN EOF",
		},
		{
			name:   "tabs between tokens",
			source: "(\t)",
			output: "LEFT_PAREN RIGHT_PAREN EOF",
		},
		{
			name:   "carriage return between tokens",
			source: "(\r)",
			output: "LEFT_PAREN RIGHT_PAREN EOF",
		},
		{
			name:   "multiple spaces",
			source: "(   )",
			output: "LEFT_PAREN RIGHT_PAREN EOF",
		},
		{
			name:   "mixed whitespace",
			source: "( \t\r )",
			output: "LEFT_PAREN RIGHT_PAREN EOF",
		},
		{
			name:   "leading and trailing whitespace",
			source: "  ( )  ",
			output: "LEFT_PAREN RIGHT_PAREN EOF",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := scanner.NewScanner(tc.source)
			got := s.ScanTokens()
			// Log tokens for debugging
			t.Logf("Source: %q", tc.source)
			for i, token := range got {
				t.Logf("[%d] %s", i, token.TokenType.String())
			}
			tokenTypes := tokenListAsString(got)
			if tokenTypes != tc.output {
				t.Fatalf("expected token types %q, got %q", tc.output, tokenTypes)
			}
		})
	}
}

func TestScanner_NewlineHandling_CorrectlyTracksLines(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		source       string
		output       string
		expectedLine int
	}{
		{
			name:         "single newline",
			source:       "(\n)",
			output:       "LEFT_PAREN RIGHT_PAREN EOF",
			expectedLine: 2, // should end on line 2
		},
		{
			name:         "multiple newlines",
			source:       "(\n\n\n)",
			output:       "LEFT_PAREN RIGHT_PAREN EOF",
			expectedLine: 4, // should end on line 4
		},
		{
			name:         "newlines with other tokens",
			source:       "{\n=\n}",
			output:       "LEFT_BRACE EQUAL RIGHT_BRACE EOF",
			expectedLine: 3, // should end on line 3
		},
		{
			name:         "no newlines",
			source:       "()",
			output:       "LEFT_PAREN RIGHT_PAREN EOF",
			expectedLine: 1, // should stay on line 1
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := scanner.NewScanner(tc.source)
			got := s.ScanTokens()
			// Log tokens for debugging
			t.Logf("Source: %q", tc.source)
			for i, token := range got {
				t.Logf("[%d] %s (line %d)", i, token.TokenType.String(), token.Line)
			}
			tokenTypes := tokenListAsString(got)
			if tokenTypes != tc.output {
				t.Fatalf("expected token types %q, got %q", tc.output, tokenTypes)
			}
			// Check that the last token (EOF) is on the expected line
			eofToken := got[len(got)-1]
			if eofToken.Line != tc.expectedLine {
				t.Fatalf("expected EOF token to be on line %d, got line %d", tc.expectedLine, eofToken.Line)
			}
		})
	}
}

func TestScanner_Numbers_CorrectlyParsed(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		source          string
		output          string
		expectErr       bool
		expectedLiteral interface{}
	}{
		{
			name:            "integer",
			source:          "123",
			output:          "NUMBER EOF",
			expectErr:       false,
			expectedLiteral: 123.0,
		},
		{
			name:            "decimal",
			source:          "123.45",
			output:          "NUMBER EOF",
			expectErr:       false,
			expectedLiteral: 123.45,
		},
		{
			name:            "zero",
			source:          "0",
			output:          "NUMBER EOF",
			expectErr:       false,
			expectedLiteral: 0.0,
		},
		{
			name:            "decimal starting with zero",
			source:          "0.123",
			output:          "NUMBER EOF",
			expectErr:       false,
			expectedLiteral: 0.123,
		},
		{
			name:            "large number",
			source:          "999999",
			output:          "NUMBER EOF",
			expectErr:       false,
			expectedLiteral: 999999.0,
		},
		{
			name:            "small decimal",
			source:          "0.001",
			output:          "NUMBER EOF",
			expectErr:       false,
			expectedLiteral: 0.001,
		},
		{
			name:            "multiple numbers",
			source:          "123 456.78",
			output:          "NUMBER NUMBER EOF",
			expectErr:       false,
			expectedLiteral: nil, // We'll check both literals separately
		},
		{
			name:            "number with operators",
			source:          "123+456",
			output:          "NUMBER PLUS NUMBER EOF",
			expectErr:       false,
			expectedLiteral: nil, // Multiple tokens
		},
		{
			name:            "number followed by dot and non-digit",
			source:          "123.abc",
			output:          "NUMBER DOT EOF", // 123, then ., then error on 'a'
			expectErr:       true,             // Will error on 'a'
			expectedLiteral: 123.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := scanner.NewScanner(tc.source)
			got := s.ScanTokens()

			// Log tokens for debugging
			t.Logf("Source: %q", tc.source)
			for i, token := range got {
				t.Logf("[%d] %s (literal: %v)", i, token.TokenType.String(), token.Literal)
			}

			if !tc.expectErr {
				tokenTypes := tokenListAsString(got)
				if tokenTypes != tc.output {
					t.Errorf("expected token types %q, got %q", tc.output, tokenTypes)
				}

				// Check specific literal for single number tests
				if tc.expectedLiteral != nil && len(got) >= 1 {
					if got[0].TokenType == token.NUMBER {
						if got[0].Literal != tc.expectedLiteral {
							t.Errorf("expected literal %v, got %v", tc.expectedLiteral, got[0].Literal)
						}
					}
				}
			}
		})
	}
}

func TestScanner_Numbers_EdgeCases(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		source    string
		output    string
		expectErr bool
	}{
		{
			name:      "standalone dot",
			source:    ".",
			output:    "DOT EOF",
			expectErr: false,
		},
		{
			name:      "dot followed by non-digit",
			source:    ".abc",
			output:    "DOT EOF", // Will error on 'a'
			expectErr: true,
		},
		{
			name:      "number expression",
			source:    "(123.45)",
			output:    "LEFT_PAREN NUMBER RIGHT_PAREN EOF",
			expectErr: false,
		},
		{
			name:      "decimal arithmetic",
			source:    "12.34 + 56.78",
			output:    "NUMBER PLUS NUMBER EOF",
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := scanner.NewScanner(tc.source)
			got := s.ScanTokens()

			// Log tokens for debugging
			t.Logf("Source: %q", tc.source)
			for i, token := range got {
				t.Logf("[%d] %s (literal: %v)", i, token.TokenType.String(), token.Literal)
			}

			if !tc.expectErr {
				tokenTypes := tokenListAsString(got)
				if tokenTypes != tc.output {
					t.Errorf("expected token types %q, got %q", tc.output, tokenTypes)
				}
			}
		})
	}
}
