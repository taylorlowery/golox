package ast

import (
	"fmt"
)

type AstPrinter struct{}

func (a *AstPrinter) visitBinaryExpr(expr *Binary) any {
	return a.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (a *AstPrinter) visitGroupingExpr(expr *Grouping) any {
	return a.parenthesize("group", expr.expression)
}

func (a *AstPrinter) visitLiteralExpr(expr *Literal) any {
	if expr.value == nil {
		return nil
	}
	return fmt.Sprint(expr.value)
}

func (a *AstPrinter) visitUnaryExpr(expr *Unary) any {
	return a.parenthesize(expr.operator.Lexeme, expr.right)
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
