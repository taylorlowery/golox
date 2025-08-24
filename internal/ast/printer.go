package ast

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type AstPrinter struct {
	Stdout io.Writer
}

type astOption func(a *AstPrinter) error

func NewAstPrinter(opts ...astOption) (*AstPrinter, error) {
	a := AstPrinter{
		Stdout: os.Stdout,
	}
	for _, opt := range opts {
		err := opt(&a)
		if err != nil {
			return nil, err
		}
	}
	return &a, nil
}

func WithStdout(w io.Writer) astOption {
	return func(a *AstPrinter) error {
		if w == nil {
			return errors.New("nil output writer")
		}
		a.Stdout = w
		return nil
	}
}

func (a *AstPrinter) visitBinaryExpr(expr *Binary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) visitGroupingExpr(expr *Grouping) any {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) visitLiteralExpr(expr *Literal) any {
	if expr.Value == nil {
		return nil
	}
	return fmt.Sprint(expr.Value)
}

func (a *AstPrinter) visitUnaryExpr(expr *Unary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) PrintAst(expr Expr) string {
	return fmt.Sprint(expr.accept(a))
}

func (a *AstPrinter) parenthesize(lexeme string, exprs ...Expr) string {
	var result string
	result += "("

	result += lexeme
	for _, e := range exprs {
		result += " "
		result += fmt.Sprint(e.accept(a))
	}

	result += ")"
	return result
}
