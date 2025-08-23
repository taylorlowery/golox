package ast

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/taylorlowery/lox/internal/token"
)

func TestPrintAst_PrintsExpected(t *testing.T) {
	t.Parallel()

	// TODO: implement test once printer is complete
	want := "(* (- 123) (group 45.67))"
	var expr = Binary{
		left: &Unary{
			operator: token.Token{
				TokenType: token.MINUS,
				Lexeme:    "-",
				Literal:   nil,
				Line:      1,
			},
			right: &Literal{
				value: 123,
			},
		},
		operator: token.Token{
			TokenType: token.STAR,
			Lexeme:    "*",
			Literal:   nil,
			Line:      1,
		},
		right: &Grouping{
			expression: &Literal{
				value: 45.67,
			},
		},
	}

	p := AstPrinter{}
	got := p.PrintAst(&expr)
	t.Log(got)
	if got != want {
		t.Fatal(cmp.Diff(want, got))
	}
}
