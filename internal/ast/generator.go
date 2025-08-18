package ast

import (
	"fmt"
	"io"
	"strings"
)

// defineType generates the code for a go struct and writes it to a given writer.
// The user supplies the struct name, and fields as a list of string.
// TODO: example usage.
func defineType(w io.Writer, structName string, fieldList string) {
	fmt.Fprintf(w, "type %s struct {\n", structName)
	for _, field := range strings.Split(fieldList, ", ") {
		fmt.Fprintf(w, "\t%s\n", field)
	}
	fmt.Fprintf(w, "}")
}
