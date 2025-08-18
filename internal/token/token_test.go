package token_test

import (
	"testing"

	"github.com/taylorlowery/lox/internal/token"
)

func TestTokenString_OutputsExpectedString(t *testing.T) {
	t.Parallel()
	token := token.Token{
		TokenType: token.AND,
		Lexeme:    "AND",
		Literal:   "AND",
		Line:      42,
	}
	got := token.String()
	want := "AND AND AND"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}
