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
	baseName := "Base"

	defineType(&output, className, baseName, fieldList)

	want := `type Example struct {
	field1 string
	field2 int
	field3 OtherType
}

func (e *Example) accept(v Visitor[any]) any {
	return v.visitExampleBase(e)
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

import "github.com/taylorlowery/lox/internal/token"

type Visitor[K any] interface {
	visitBinaryExpr(b *Binary) K
	visitGroupingExpr(g *Grouping) K
	visitLiteralExpr(l *Literal) K
	visitUnaryExpr(u *Unary) K
}

type Expr interface{
	accept(v Visitor[any]) any
}

type Binary struct {
	left Expr
	operator token.Token
	right Expr
}

func (b *Binary) accept(v Visitor[any]) any {
	return v.visitBinaryExpr(b)
}


type Grouping struct {
	expression Expr
}

func (g *Grouping) accept(v Visitor[any]) any {
	return v.visitGroupingExpr(g)
}


type Literal struct {
	value any
}

func (l *Literal) accept(v Visitor[any]) any {
	return v.visitLiteralExpr(l)
}


type Unary struct {
	operator token.Token
	right Expr
}

func (u *Unary) accept(v Visitor[any]) any {
	return v.visitUnaryExpr(u)
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

	defineVisitor(&output, "Expr", typeDefs)

	want := `type Visitor[K any] interface {
	visitBinaryExpr(b *Binary) K
	visitGroupingExpr(g *Grouping) K
	visitLiteralExpr(l *Literal) K
	visitUnaryExpr(u *Unary) K
}

`
	got := output.String()
	t.Log(got)

	if got != want {
		t.Fatal(cmp.Diff(want, got))
	}
}
