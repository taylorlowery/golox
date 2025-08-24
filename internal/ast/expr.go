package ast

import "github.com/taylorlowery/lox/internal/token"

type Visitor[K any] interface {
	visitBinaryExpr(b *Binary) K
	visitGroupingExpr(g *Grouping) K
	visitLiteralExpr(l *Literal) K
	visitUnaryExpr(u *Unary) K
}

type Expr interface {
	accept(v Visitor[any]) any
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b *Binary) accept(v Visitor[any]) any {
	return v.visitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) accept(v Visitor[any]) any {
	return v.visitGroupingExpr(g)
}

type Literal struct {
	Value any
}

func (l *Literal) accept(v Visitor[any]) any {
	return v.visitLiteralExpr(l)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u *Unary) accept(v Visitor[any]) any {
	return v.visitUnaryExpr(u)
}
