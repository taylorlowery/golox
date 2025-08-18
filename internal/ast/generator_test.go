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
}`
	got := output.String()

	if got != want {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestDefineAst_GeneratesCodeFileWithAllExpectedStructs(t *testing.T) {
	t.Parallel()

	// var output bytes.Buffer

}
