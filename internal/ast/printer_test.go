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
		Left: &Unary{
			Operator: token.Token{
				TokenType: token.MINUS,
				Lexeme:    "-",
				Literal:   nil,
				Line:      1,
			},
			Right: &Literal{
				Value: 123,
			},
		},
		Operator: token.Token{
			TokenType: token.STAR,
			Lexeme:    "*",
			Literal:   nil,
			Line:      1,
		},
		Right: &Grouping{
			Expression: &Literal{
				Value: 45.67,
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
