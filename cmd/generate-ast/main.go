package main

import (
	"fmt"
	"os"

	"github.com/taylorlowery/lox/internal/ast"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: generate-ast <output directory>")
		os.Exit(64)
	}

	outputFile := os.Args[1]

	packageName := "golox"

	typeDefs := []string{
		"Binary   : left Expr, operator token.Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value any",
		"Unary    : operator token.Token, right Expr",
	}

	err := ast.GenerateAst(outputFile, packageName, typeDefs)
	if err != nil {
		fmt.Println(err)
		os.Exit(64)
	}

}
