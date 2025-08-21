package golox

import "github.com/taylorlowery/lox/internal/token"

type Visitor[K any] interface {
	visitBinary(b Binary) K
	visitGrouping(g Grouping) K
	visitLiteral(l Literal) K
	visitUnary(u Unary) K
}

type Expr interface{
	accept(v Visitor[any]) any
}

type Binary struct {
	left Expr
	operator token.Token
	right Expr
}

func (b Binary) accept(v Visitor[any]) any {
	return v.visitBinary(b)
}


type Grouping struct {
	expression Expr
}

func (g Grouping) accept(v Visitor[any]) any {
	return v.visitGrouping(g)
}


type Literal struct {
	value any
}

func (l Literal) accept(v Visitor[any]) any {
	return v.visitLiteral(l)
}


type Unary struct {
	operator token.Token
	right Expr
}

func (u Unary) accept(v Visitor[any]) any {
	return v.visitUnary(u)
}


