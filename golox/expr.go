package golox

import (
	"github.com/taylorlowery/lox/internal/token"
)

type Expr interface{}

type Binary struct {
	left     Expr
	operator token.Token
	right    Expr
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value any
}

type Unary struct {
	operator token.Token
	right    Expr
}
