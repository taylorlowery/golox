package ast

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDefineType_OutputsExpectedCode(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer

	className := "Example"
	fieldList := "field1 string, field2 int, field3 OtherType"

	defineType(&output, className, fieldList)

	want := `type Example struct {
	field1 string
	field2 int
	field3 OtherType
}`
	got := output.String()

	if got != want {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestDefineAst_GeneratesCodeFileWithAllExpectedStructs(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer

	typesDefs := []string{
		"Binary   : left Expr, operator token.Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value any",
		"Unary    : operator token.Token, right Expr",
	}

	defineAst(&output, "golox", "Expr", typesDefs)

	want := `package golox

type Expr interface{}

type Binary struct {
	left Expr
	operator token.Token
	right Expr
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value any
}

type Unary struct {
	operator token.Token
	right Expr
}

`
	got := output.String()

	if got != want {
		t.Fatal(cmp.Diff(want, got))
	}

}
