package parser

import (
	"reflect"
	"testing"

	"github.com/taylorlowery/lox/internal/ast"
	"github.com/taylorlowery/lox/internal/token"
)

// Helper function to create tokens easily
func makeToken(tokenType token.TokenType, lexeme string, literal any) token.Token {
	return token.Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      1,
	}
}

// Helper function to create a simple token list with EOF
func makeTokens(tokens ...token.Token) []token.Token {
	result := make([]token.Token, len(tokens)+1)
	copy(result, tokens)
	result[len(tokens)] = makeToken(token.EOF, "", nil)
	return result
}

func TestNewParser(t *testing.T) {
	t.Parallel()

	tokens := []token.Token{
		makeToken(token.NUMBER, "42", 42.0),
		makeToken(token.EOF, "", nil),
	}

	parser := NewParser(tokens)

	if parser == nil {
		t.Fatal("NewParser returned nil")
	}
	if !reflect.DeepEqual(parser.tokens, tokens) {
		t.Errorf("Expected tokens %v, got %v", tokens, parser.tokens)
	}
	if parser.current != 0 {
		t.Errorf("Expected current to be 0, got %d", parser.current)
	}
}

func TestParser_PrimaryExpressions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		tokens   []token.Token
		expected ast.Expr
	}{
		{
			name:   "number literal",
			tokens: makeTokens(makeToken(token.NUMBER, "42", 42.0)),
			expected: &ast.Literal{
				Value: 42.0,
			},
		},
		{
			name:   "string literal",
			tokens: makeTokens(makeToken(token.STRING, "\"hello\"", "hello")),
			expected: &ast.Literal{
				Value: "hello",
			},
		},
		{
			name:   "true literal",
			tokens: makeTokens(makeToken(token.TRUE, "true", nil)),
			expected: &ast.Literal{
				Value: true,
			},
		},
		{
			name:   "false literal",
			tokens: makeTokens(makeToken(token.FALSE, "false", nil)),
			expected: &ast.Literal{
				Value: false,
			},
		},
		{
			name:   "nil literal",
			tokens: makeTokens(makeToken(token.NIL, "nil", nil)),
			expected: &ast.Literal{
				Value: nil,
			},
		},
		{
			name: "grouped expression",
			tokens: makeTokens(
				makeToken(token.LEFT_PAREN, "(", nil),
				makeToken(token.NUMBER, "42", 42.0),
				makeToken(token.RIGHT_PAREN, ")", nil),
			),
			expected: &ast.Grouping{
				Expression: &ast.Literal{
					Value: 42.0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(tc.tokens)
			result := parser.expression()

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestParser_UnaryExpressions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		tokens   []token.Token
		expected ast.Expr
	}{
		{
			name: "negation",
			tokens: makeTokens(
				makeToken(token.MINUS, "-", nil),
				makeToken(token.NUMBER, "42", 42.0),
			),
			expected: &ast.Unary{
				Operator: makeToken(token.MINUS, "-", nil),
				Right: &ast.Literal{
					Value: 42.0,
				},
			},
		},
		{
			name: "logical not",
			tokens: makeTokens(
				makeToken(token.BANG, "!", nil),
				makeToken(token.TRUE, "true", nil),
			),
			expected: &ast.Unary{
				Operator: makeToken(token.BANG, "!", nil),
				Right: &ast.Literal{
					Value: true,
				},
			},
		},
		{
			name: "double negation",
			tokens: makeTokens(
				makeToken(token.MINUS, "-", nil),
				makeToken(token.MINUS, "-", nil),
				makeToken(token.NUMBER, "42", 42.0),
			),
			expected: &ast.Unary{
				Operator: makeToken(token.MINUS, "-", nil),
				Right: &ast.Unary{
					Operator: makeToken(token.MINUS, "-", nil),
					Right: &ast.Literal{
						Value: 42.0,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(tc.tokens)
			result := parser.expression()

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestParser_BinaryExpressions_Factor(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		tokens   []token.Token
		expected ast.Expr
	}{
		{
			name: "multiplication",
			tokens: makeTokens(
				makeToken(token.NUMBER, "6", 6.0),
				makeToken(token.STAR, "*", nil),
				makeToken(token.NUMBER, "7", 7.0),
			),
			expected: &ast.Binary{
				Left: &ast.Literal{
					Value: 6.0,
				},
				Operator: makeToken(token.STAR, "*", nil),
				Right: &ast.Literal{
					Value: 7.0,
				},
			},
		},
		{
			name: "division",
			tokens: makeTokens(
				makeToken(token.NUMBER, "10", 10.0),
				makeToken(token.SLASH, "/", nil),
				makeToken(token.NUMBER, "2", 2.0),
			),
			expected: &ast.Binary{
				Left: &ast.Literal{
					Value: 10.0,
				},
				Operator: makeToken(token.SLASH, "/", nil),
				Right: &ast.Literal{
					Value: 2.0,
				},
			},
		},
		{
			name: "left associative multiplication",
			tokens: makeTokens(
				makeToken(token.NUMBER, "2", 2.0),
				makeToken(token.STAR, "*", nil),
				makeToken(token.NUMBER, "3", 3.0),
				makeToken(token.STAR, "*", nil),
				makeToken(token.NUMBER, "4", 4.0),
			),
			expected: &ast.Binary{
				Left: &ast.Binary{
					Left: &ast.Literal{
						Value: 2.0,
					},
					Operator: makeToken(token.STAR, "*", nil),
					Right: &ast.Literal{
						Value: 3.0,
					},
				},
				Operator: makeToken(token.STAR, "*", nil),
				Right: &ast.Literal{
					Value: 4.0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(tc.tokens)
			result := parser.expression()

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestParser_BinaryExpressions_Term(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		tokens   []token.Token
		expected ast.Expr
	}{
		{
			name: "addition",
			tokens: makeTokens(
				makeToken(token.NUMBER, "1", 1.0),
				makeToken(token.PLUS, "+", nil),
				makeToken(token.NUMBER, "2", 2.0),
			),
			expected: &ast.Binary{
				Left: &ast.Literal{
					Value: 1.0,
				},
				Operator: makeToken(token.PLUS, "+", nil),
				Right: &ast.Literal{
					Value: 2.0,
				},
			},
		},
		{
			name: "subtraction",
			tokens: makeTokens(
				makeToken(token.NUMBER, "5", 5.0),
				makeToken(token.MINUS, "-", nil),
				makeToken(token.NUMBER, "3", 3.0),
			),
			expected: &ast.Binary{
				Left: &ast.Literal{
					Value: 5.0,
				},
				Operator: makeToken(token.MINUS, "-", nil),
				Right: &ast.Literal{
					Value: 3.0,
				},
			},
		},
		{
			name: "left associative addition",
			tokens: makeTokens(
				makeToken(token.NUMBER, "1", 1.0),
				makeToken(token.PLUS, "+", nil),
				makeToken(token.NUMBER, "2", 2.0),
				makeToken(token.PLUS, "+", nil),
				makeToken(token.NUMBER, "3", 3.0),
			),
			expected: &ast.Binary{
				Left: &ast.Binary{
					Left: &ast.Literal{
						Value: 1.0,
					},
					Operator: makeToken(token.PLUS, "+", nil),
					Right: &ast.Literal{
						Value: 2.0,
					},
				},
				Operator: makeToken(token.PLUS, "+", nil),
				Right: &ast.Literal{
					Value: 3.0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(tc.tokens)
			result := parser.expression()

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestParser_BinaryExpressions_Comparison(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		tokens   []token.Token
		expected ast.Expr
	}{
		{
			name: "greater than",
			tokens: makeTokens(
				makeToken(token.NUMBER, "5", 5.0),
				makeToken(token.GREATER, ">", nil),
				makeToken(token.NUMBER, "3", 3.0),
			),
			expected: &ast.Binary{
				Left: &ast.Literal{
					Value: 5.0,
				},
				Operator: makeToken(token.GREATER, ">", nil),
				Right: &ast.Literal{
					Value: 3.0,
				},
			},
		},
		{
			name: "greater than or equal",
			tokens: makeTokens(
				makeToken(token.NUMBER, "5", 5.0),
				makeToken(token.GREATER_EQUAL, ">=", nil),
				makeToken(token.NUMBER, "5", 5.0),
			),
			expected: &ast.Binary{
				Left: &ast.Literal{
					Value: 5.0,
				},
				Operator: makeToken(token.GREATER_EQUAL, ">=", nil),
				Right: &ast.Literal{
					Value: 5.0,
				},
			},
		},
		{
			name: "less than",
			tokens: makeTokens(
				makeToken(token.NUMBER, "3", 3.0),
				makeToken(token.LESS, "<", nil),
				makeToken(token.NUMBER, "5", 5.0),
			),
			expected: &ast.Binary{
				Left: &ast.Literal{
					Value: 3.0,
				},
				Operator: makeToken(token.LESS, "<", nil),
				Right: &ast.Literal{
					Value: 5.0,
				},
			},
		},
		{
			name: "less than or equal",
			tokens: makeTokens(
				makeToken(token.NUMBER, "3", 3.0),
				makeToken(token.LESS_EQUAL, "<=", nil),
				makeToken(token.NUMBER, "5", 5.0),
			),
			expected: &ast.Binary{
				Left: &ast.Literal{
					Value: 3.0,
				},
				Operator: makeToken(token.LESS_EQUAL, "<=", nil),
				Right: &ast.Literal{
					Value: 5.0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(tc.tokens)
			result := parser.expression()

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestParser_BinaryExpressions_Equality(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		tokens   []token.Token
		expected ast.Expr
	}{
		{
			name: "equality",
			tokens: makeTokens(
				makeToken(token.NUMBER, "5", 5.0),
				makeToken(token.EQUAL_EQUAL, "==", nil),
				makeToken(token.NUMBER, "5", 5.0),
			),
			expected: &ast.Binary{
				Left: &ast.Literal{
					Value: 5.0,
				},
				Operator: makeToken(token.EQUAL_EQUAL, "==", nil),
				Right: &ast.Literal{
					Value: 5.0,
				},
			},
		},
		{
			name: "inequality",
			tokens: makeTokens(
				makeToken(token.NUMBER, "5", 5.0),
				makeToken(token.BANG_EQUAL, "!=", nil),
				makeToken(token.NUMBER, "3", 3.0),
			),
			expected: &ast.Binary{
				Left: &ast.Literal{
					Value: 5.0,
				},
				Operator: makeToken(token.BANG_EQUAL, "!=", nil),
				Right: &ast.Literal{
					Value: 3.0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(tc.tokens)
			result := parser.expression()

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestParser_OperatorPrecedence(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		tokens   []token.Token
		expected ast.Expr
	}{
		{
			name: "multiplication before addition",
			tokens: makeTokens(
				makeToken(token.NUMBER, "2", 2.0),
				makeToken(token.PLUS, "+", nil),
				makeToken(token.NUMBER, "3", 3.0),
				makeToken(token.STAR, "*", nil),
				makeToken(token.NUMBER, "4", 4.0),
			),
			expected: &ast.Binary{
				Left: &ast.Literal{
					Value: 2.0,
				},
				Operator: makeToken(token.PLUS, "+", nil),
				Right: &ast.Binary{
					Left: &ast.Literal{
						Value: 3.0,
					},
					Operator: makeToken(token.STAR, "*", nil),
					Right: &ast.Literal{
						Value: 4.0,
					},
				},
			},
		},
		{
			name: "unary before multiplication",
			tokens: makeTokens(
				makeToken(token.MINUS, "-", nil),
				makeToken(token.NUMBER, "2", 2.0),
				makeToken(token.STAR, "*", nil),
				makeToken(token.NUMBER, "3", 3.0),
			),
			expected: &ast.Binary{
				Left: &ast.Unary{
					Operator: makeToken(token.MINUS, "-", nil),
					Right: &ast.Literal{
						Value: 2.0,
					},
				},
				Operator: makeToken(token.STAR, "*", nil),
				Right: &ast.Literal{
					Value: 3.0,
				},
			},
		},
		{
			name: "comparison before equality",
			tokens: makeTokens(
				makeToken(token.NUMBER, "1", 1.0),
				makeToken(token.LESS, "<", nil),
				makeToken(token.NUMBER, "2", 2.0),
				makeToken(token.EQUAL_EQUAL, "==", nil),
				makeToken(token.TRUE, "true", nil),
			),
			expected: &ast.Binary{
				Left: &ast.Binary{
					Left: &ast.Literal{
						Value: 1.0,
					},
					Operator: makeToken(token.LESS, "<", nil),
					Right: &ast.Literal{
						Value: 2.0,
					},
				},
				Operator: makeToken(token.EQUAL_EQUAL, "==", nil),
				Right: &ast.Literal{
					Value: true,
				},
			},
		},
		{
			name: "complex precedence with parentheses",
			tokens: makeTokens(
				makeToken(token.LEFT_PAREN, "(", nil),
				makeToken(token.NUMBER, "1", 1.0),
				makeToken(token.PLUS, "+", nil),
				makeToken(token.NUMBER, "2", 2.0),
				makeToken(token.RIGHT_PAREN, ")", nil),
				makeToken(token.STAR, "*", nil),
				makeToken(token.NUMBER, "3", 3.0),
			),
			expected: &ast.Binary{
				Left: &ast.Grouping{
					Expression: &ast.Binary{
						Left: &ast.Literal{
							Value: 1.0,
						},
						Operator: makeToken(token.PLUS, "+", nil),
						Right: &ast.Literal{
							Value: 2.0,
						},
					},
				},
				Operator: makeToken(token.STAR, "*", nil),
				Right: &ast.Literal{
					Value: 3.0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(tc.tokens)
			result := parser.expression()

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestParser_HelperMethods(t *testing.T) {
	t.Parallel()

	t.Run("peek", func(t *testing.T) {
		tokens := makeTokens(
			makeToken(token.NUMBER, "42", 42.0),
			makeToken(token.PLUS, "+", nil),
		)
		parser := NewParser(tokens)

		// Should return first token without advancing
		result := parser.peek()
		expected := makeToken(token.NUMBER, "42", 42.0)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
		if parser.current != 0 {
			t.Errorf("Expected current to remain 0, got %d", parser.current)
		}
	})

	t.Run("advance", func(t *testing.T) {
		tokens := makeTokens(
			makeToken(token.NUMBER, "42", 42.0),
			makeToken(token.PLUS, "+", nil),
		)
		parser := NewParser(tokens)

		// Should return current token and advance
		result := parser.advance()
		expected := makeToken(token.NUMBER, "42", 42.0)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
		if parser.current != 1 {
			t.Errorf("Expected current to be 1, got %d", parser.current)
		}
	})

	t.Run("previous", func(t *testing.T) {
		tokens := makeTokens(
			makeToken(token.NUMBER, "42", 42.0),
			makeToken(token.PLUS, "+", nil),
		)
		parser := NewParser(tokens)
		parser.advance() // Move to position 1

		result := parser.previous()
		expected := makeToken(token.NUMBER, "42", 42.0)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("isAtEnd false", func(t *testing.T) {
		tokens := makeTokens(makeToken(token.NUMBER, "42", 42.0))
		parser := NewParser(tokens)

		if parser.isAtEnd() {
			t.Error("Expected isAtEnd to be false")
		}
	})

	t.Run("isAtEnd true", func(t *testing.T) {
		tokens := makeTokens(makeToken(token.NUMBER, "42", 42.0))
		parser := NewParser(tokens)
		parser.advance() // Move to EOF

		if !parser.isAtEnd() {
			t.Error("Expected isAtEnd to be true")
		}
	})

	t.Run("check true", func(t *testing.T) {
		tokens := makeTokens(makeToken(token.NUMBER, "42", 42.0))
		parser := NewParser(tokens)

		if !parser.check(token.NUMBER) {
			t.Error("Expected check to return true for NUMBER token")
		}
	})

	t.Run("check false", func(t *testing.T) {
		tokens := makeTokens(makeToken(token.NUMBER, "42", 42.0))
		parser := NewParser(tokens)

		if parser.check(token.STRING) {
			t.Error("Expected check to return false for STRING token")
		}
	})

	t.Run("match true", func(t *testing.T) {
		tokens := makeTokens(makeToken(token.NUMBER, "42", 42.0))
		parser := NewParser(tokens)

		if !parser.match(token.NUMBER, token.STRING) {
			t.Error("Expected match to return true")
		}
		if parser.current != 1 {
			t.Errorf("Expected current to be 1 after match, got %d", parser.current)
		}
	})

	t.Run("match false", func(t *testing.T) {
		tokens := makeTokens(makeToken(token.NUMBER, "42", 42.0))
		parser := NewParser(tokens)

		if parser.match(token.STRING, token.TRUE) {
			t.Error("Expected match to return false")
		}
		if parser.current != 0 {
			t.Errorf("Expected current to remain 0 after failed match, got %d", parser.current)
		}
	})
}

func TestParser_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("missing right parenthesis", func(t *testing.T) {
		tokens := makeTokens(
			makeToken(token.LEFT_PAREN, "(", nil),
			makeToken(token.NUMBER, "42", 42.0),
			// Missing RIGHT_PAREN
		)
		parser := NewParser(tokens)

		// Should panic with parse error
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for missing right parenthesis")
			}
		}()
		parser.expression()
	})

	t.Run("parseError", func(t *testing.T) {
		tokens := makeTokens(makeToken(token.NUMBER, "42", 42.0))
		parser := NewParser(tokens)

		tok := makeToken(token.NUMBER, "42", 42.0)
		err := parser.parseError(tok, "test message")

		if err == nil {
			t.Error("Expected parseError to return an error")
		}
		if err.Error() != "Parser error: test message" {
			t.Errorf("Expected 'Parser error: test message', got %s", err.Error())
		}
	})

	t.Run("consume success", func(t *testing.T) {
		tokens := makeTokens(
			makeToken(token.LEFT_PAREN, "(", nil),
			makeToken(token.RIGHT_PAREN, ")", nil),
		)
		parser := NewParser(tokens)

		result := parser.consume(token.LEFT_PAREN, "expected (")
		expected := makeToken(token.LEFT_PAREN, "(", nil)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("consume failure", func(t *testing.T) {
		tokens := makeTokens(makeToken(token.NUMBER, "42", 42.0))
		parser := NewParser(tokens)

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic from failed consume")
			}
		}()
		parser.consume(token.LEFT_PAREN, "expected (")
	})
}

func TestParser_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty tokens returns nil", func(t *testing.T) {
		tokens := []token.Token{makeToken(token.EOF, "", nil)}
		parser := NewParser(tokens)

		result := parser.expression()
		if result != nil {
			t.Errorf("Expected nil for empty input, got %+v", result)
		}
	})

	t.Run("nested groupings", func(t *testing.T) {
		tokens := makeTokens(
			makeToken(token.LEFT_PAREN, "(", nil),
			makeToken(token.LEFT_PAREN, "(", nil),
			makeToken(token.NUMBER, "42", 42.0),
			makeToken(token.RIGHT_PAREN, ")", nil),
			makeToken(token.RIGHT_PAREN, ")", nil),
		)
		parser := NewParser(tokens)

		result := parser.expression()
		expected := &ast.Grouping{
			Expression: &ast.Grouping{
				Expression: &ast.Literal{
					Value: 42.0,
				},
			},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("advance at end doesn't go past EOF", func(t *testing.T) {
		tokens := makeTokens(makeToken(token.NUMBER, "42", 42.0))
		parser := NewParser(tokens)
		parser.advance() // Move to EOF position

		// Advancing at EOF should return the previous token and not advance current
		result := parser.advance()
		expected := makeToken(token.NUMBER, "42", 42.0)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
		if parser.current != 1 {
			t.Errorf("Expected current to remain 1 at EOF, got %d", parser.current)
		}
	})
}

func TestParser_ComplexExpressions(t *testing.T) {
	t.Parallel()

	t.Run("mixed arithmetic with precedence", func(t *testing.T) {
		// 1 + 2 * 3 - 4 / 2
		tokens := makeTokens(
			makeToken(token.NUMBER, "1", 1.0),
			makeToken(token.PLUS, "+", nil),
			makeToken(token.NUMBER, "2", 2.0),
			makeToken(token.STAR, "*", nil),
			makeToken(token.NUMBER, "3", 3.0),
			makeToken(token.MINUS, "-", nil),
			makeToken(token.NUMBER, "4", 4.0),
			makeToken(token.SLASH, "/", nil),
			makeToken(token.NUMBER, "2", 2.0),
		)
		parser := NewParser(tokens)

		result := parser.expression()

		// Should parse as: ((1 + (2 * 3)) - (4 / 2))
		expected := &ast.Binary{
			Left: &ast.Binary{
				Left: &ast.Literal{
					Value: 1.0,
				},
				Operator: makeToken(token.PLUS, "+", nil),
				Right: &ast.Binary{
					Left: &ast.Literal{
						Value: 2.0,
					},
					Operator: makeToken(token.STAR, "*", nil),
					Right: &ast.Literal{
						Value: 3.0,
					},
				},
			},
			Operator: makeToken(token.MINUS, "-", nil),
			Right: &ast.Binary{
				Left: &ast.Literal{
					Value: 4.0,
				},
				Operator: makeToken(token.SLASH, "/", nil),
				Right: &ast.Literal{
					Value: 2.0,
				},
			},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("complex comparison and equality", func(t *testing.T) {
		// 1 + 2 > 3 == false
		tokens := makeTokens(
			makeToken(token.NUMBER, "1", 1.0),
			makeToken(token.PLUS, "+", nil),
			makeToken(token.NUMBER, "2", 2.0),
			makeToken(token.GREATER, ">", nil),
			makeToken(token.NUMBER, "3", 3.0),
			makeToken(token.EQUAL_EQUAL, "==", nil),
			makeToken(token.FALSE, "false", nil),
		)
		parser := NewParser(tokens)

		result := parser.expression()

		// Should parse as: (((1 + 2) > 3) == false)
		expected := &ast.Binary{
			Left: &ast.Binary{
				Left: &ast.Binary{
					Left: &ast.Literal{
						Value: 1.0,
					},
					Operator: makeToken(token.PLUS, "+", nil),
					Right: &ast.Literal{
						Value: 2.0,
					},
				},
				Operator: makeToken(token.GREATER, ">", nil),
				Right: &ast.Literal{
					Value: 3.0,
				},
			},
			Operator: makeToken(token.EQUAL_EQUAL, "==", nil),
			Right: &ast.Literal{
				Value: false,
			},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})
}

func TestParser_Synchronize(t *testing.T) {
	t.Parallel()

	t.Run("synchronize after semicolon", func(t *testing.T) {
		tokens := makeTokens(
			makeToken(token.NUMBER, "42", 42.0),
			makeToken(token.SEMICOLON, ";", nil),
			makeToken(token.STRING, "\"hello\"", "hello"),
		)
		parser := NewParser(tokens)

		// Call synchronize and verify it advances past semicolon
		parser.synchronize()

		// Should be at the string token now
		current := parser.peek()
		expected := makeToken(token.STRING, "\"hello\"", "hello")
		if !reflect.DeepEqual(current, expected) {
			t.Errorf("Expected %+v, got %+v", expected, current)
		}
	})

	t.Run("synchronize at statement keyword", func(t *testing.T) {
		tokens := makeTokens(
			makeToken(token.NUMBER, "42", 42.0),
			makeToken(token.PLUS, "+", nil),
			makeToken(token.CLASS, "class", nil),
			makeToken(token.IDENTIFIER, "MyClass", nil),
		)
		parser := NewParser(tokens)

		// Advance past first token, then synchronize
		parser.advance()
		parser.synchronize()

		// Should be at the CLASS token
		current := parser.peek()
		expected := makeToken(token.CLASS, "class", nil)
		if !reflect.DeepEqual(current, expected) {
			t.Errorf("Expected %+v, got %+v", expected, current)
		}
	})

	t.Run("synchronize at EOF", func(t *testing.T) {
		tokens := makeTokens(makeToken(token.NUMBER, "42", 42.0))
		parser := NewParser(tokens)

		// Advance to EOF then synchronize
		parser.advance()
		parser.synchronize()

		// Should still be at EOF
		if !parser.isAtEnd() {
			t.Error("Expected to be at EOF after synchronize")
		}
	})
}

func TestParser_AdditionalEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("multiple unary operators", func(t *testing.T) {
		tokens := makeTokens(
			makeToken(token.BANG, "!", nil),
			makeToken(token.BANG, "!", nil),
			makeToken(token.MINUS, "-", nil),
			makeToken(token.NUMBER, "42", 42.0),
		)
		parser := NewParser(tokens)
		result := parser.expression()

		// Should parse as: !(!(-42))
		expected := &ast.Unary{
			Operator: makeToken(token.BANG, "!", nil),
			Right: &ast.Unary{
				Operator: makeToken(token.BANG, "!", nil),
				Right: &ast.Unary{
					Operator: makeToken(token.MINUS, "-", nil),
					Right: &ast.Literal{
						Value: 42.0,
					},
				},
			},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("deeply nested parentheses", func(t *testing.T) {
		tokens := makeTokens(
			makeToken(token.LEFT_PAREN, "(", nil),
			makeToken(token.LEFT_PAREN, "(", nil),
			makeToken(token.LEFT_PAREN, "(", nil),
			makeToken(token.NUMBER, "42", 42.0),
			makeToken(token.RIGHT_PAREN, ")", nil),
			makeToken(token.RIGHT_PAREN, ")", nil),
			makeToken(token.RIGHT_PAREN, ")", nil),
		)
		parser := NewParser(tokens)
		result := parser.expression()

		expected := &ast.Grouping{
			Expression: &ast.Grouping{
				Expression: &ast.Grouping{
					Expression: &ast.Literal{
						Value: 42.0,
					},
				},
			},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("mixed types in binary expression", func(t *testing.T) {
		tokens := makeTokens(
			makeToken(token.STRING, "\"hello\"", "hello"),
			makeToken(token.PLUS, "+", nil),
			makeToken(token.NUMBER, "42", 42.0),
		)
		parser := NewParser(tokens)
		result := parser.expression()

		expected := &ast.Binary{
			Left: &ast.Literal{
				Value: "hello",
			},
			Operator: makeToken(token.PLUS, "+", nil),
			Right: &ast.Literal{
				Value: 42.0,
			},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("all boolean literals", func(t *testing.T) {
		testCases := []struct {
			name     string
			token    token.Token
			expected any
		}{
			{"true", makeToken(token.TRUE, "true", nil), true},
			{"false", makeToken(token.FALSE, "false", nil), false},
			{"nil", makeToken(token.NIL, "nil", nil), nil},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tokens := makeTokens(tc.token)
				parser := NewParser(tokens)
				result := parser.expression()

				expected := &ast.Literal{
					Value: tc.expected,
				}

				if !reflect.DeepEqual(result, expected) {
					t.Errorf("Expected %+v, got %+v", expected, result)
				}
			})
		}
	})

	t.Run("all comparison operators", func(t *testing.T) {
		operators := []token.TokenType{
			token.GREATER,
			token.GREATER_EQUAL,
			token.LESS,
			token.LESS_EQUAL,
		}

		for _, op := range operators {
			t.Run(op.String(), func(t *testing.T) {
				tokens := makeTokens(
					makeToken(token.NUMBER, "5", 5.0),
					makeToken(op, op.String(), nil),
					makeToken(token.NUMBER, "3", 3.0),
				)
				parser := NewParser(tokens)
				result := parser.expression()

				expected := &ast.Binary{
					Left: &ast.Literal{
						Value: 5.0,
					},
					Operator: makeToken(op, op.String(), nil),
					Right: &ast.Literal{
						Value: 3.0,
					},
				}

				if !reflect.DeepEqual(result, expected) {
					t.Errorf("Expected %+v, got %+v", expected, result)
				}
			})
		}
	})

	t.Run("chain of same precedence operators", func(t *testing.T) {
		// Test: 1 == 2 != 3 == 4
		tokens := makeTokens(
			makeToken(token.NUMBER, "1", 1.0),
			makeToken(token.EQUAL_EQUAL, "==", nil),
			makeToken(token.NUMBER, "2", 2.0),
			makeToken(token.BANG_EQUAL, "!=", nil),
			makeToken(token.NUMBER, "3", 3.0),
			makeToken(token.EQUAL_EQUAL, "==", nil),
			makeToken(token.NUMBER, "4", 4.0),
		)
		parser := NewParser(tokens)
		result := parser.expression()

		// Should parse as: (((1 == 2) != 3) == 4)
		expected := &ast.Binary{
			Left: &ast.Binary{
				Left: &ast.Binary{
					Left: &ast.Literal{
						Value: 1.0,
					},
					Operator: makeToken(token.EQUAL_EQUAL, "==", nil),
					Right: &ast.Literal{
						Value: 2.0,
					},
				},
				Operator: makeToken(token.BANG_EQUAL, "!=", nil),
				Right: &ast.Literal{
					Value: 3.0,
				},
			},
			Operator: makeToken(token.EQUAL_EQUAL, "==", nil),
			Right: &ast.Literal{
				Value: 4.0,
			},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})
}
