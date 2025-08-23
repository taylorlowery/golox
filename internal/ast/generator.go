package ast

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// defineType generates the code for a go struct and writes it to a given writer.
// The user supplies the struct name, and fields as a list of string.
// TODO: example usage.
func defineType(w io.Writer, structName string, baseName string, fieldList string) {
	fmt.Fprintf(w, "type %s struct {\n", structName)
	for field := range strings.SplitSeq(fieldList, ", ") {
		fmt.Fprintf(w, "\t%s\n", field)
	}
	fmt.Fprintf(w, "}\n\n")
	fmt.Fprintf(w, "func (%c *%s) accept(v Visitor[any]) any {\n\treturn v.visit%s%s(%c)\n}\n", strings.ToLower(structName)[0], structName, structName, baseName, strings.ToLower(structName)[0])
}

func defineAst(w io.Writer, packageName string, baseName string, typeDefs []string) {
	fmt.Fprintf(w, "package %s\n\n", packageName)

	fmt.Fprintf(w, "import \"github.com/taylorlowery/lox/internal/token\"\n\n")

	defineVisitor(w, baseName, typeDefs)

	fmt.Fprintf(w, "type %s interface{\n\taccept(v Visitor[any]) any\n}\n\n", baseName)

	for _, typeDef := range typeDefs {
		parts := strings.Split(typeDef, ":")
		structName := strings.TrimSpace(parts[0])
		fields := strings.TrimSpace(parts[1])
		defineType(w, structName, baseName, fields)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
	}
}

func defineVisitor(w io.Writer, baseName string, typeDefs []string) {
	fmt.Fprintf(w, "type Visitor[K any] interface {\n")

	for _, typeDef := range typeDefs {
		typeName := strings.TrimSpace(strings.Split(typeDef, ":")[0])
		fmt.Fprintf(w, "\tvisit%s%s(%c *%s) K\n", typeName, baseName, strings.ToLower(typeName)[0], typeName)
	}

	fmt.Fprintf(w, "}\n\n")
}

func GenerateAst(outputPath string, packageName string, typeDefs []string) error {
	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	defineAst(file, packageName, "Expr", typeDefs)
	return nil
}
