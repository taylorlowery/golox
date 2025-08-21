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
}

func (e *Example)[K any] accept(v Visitor[K]) K {
	return v.visitExample(e)
}

`
	got := output.String()

	if got != want {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestDefineAst_GeneratesCodeFileWithAllExpectedStructs(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer

	typeDefs := []string{
		"Binary   : left Expr, operator token.Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value any",
		"Unary    : operator token.Token, right Expr",
	}

	defineAst(&output, "golox", "Expr", typeDefs)

	want := `package golox

type Visitor[K any] interface {
	visitBinary(b Binary) K
	visitGrouping(g Grouping) K
	visitLiteral(l Literal) K
	visitUnary(u Unary) K
}

type Expr[K any] interface{
	accept(v Visitor) K
}

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

func TestDefineVisitor(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer

	typeDefs := []string{
		"Binary   : left Expr, operator token.Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value any",
		"Unary    : operator token.Token, right Expr",
	}

	defineVisitor(&output, typeDefs)

	want := `type Visitor[K any] interface {
	visitBinary(b Binary) K
	visitGrouping(g Grouping) K
	visitLiteral(l Literal) K
	visitUnary(u Unary) K
}

`
	got := output.String()
	t.Log(got)

	if got != want {
		t.Fatal(cmp.Diff(want, got))
	}
}
