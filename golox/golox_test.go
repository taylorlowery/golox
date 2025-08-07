package golox_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/taylorlowery/lox/golox"
)

func TestRunPrompt_ReadsInputAndPrintsIt(t *testing.T) {
	t.Parallel()
	input := strings.NewReader("1\n2\n3\n")
	var output bytes.Buffer
	var errOutput bytes.Buffer

	g, err := golox.NewGolox(
		golox.WithInput(input),
		golox.WithOutput(&output),
		golox.WithStderr(&errOutput),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = g.RunPrompt()
	if err != nil {
		t.Fatal("expected an error or something")
	}
	got := output.String()
	want := "> 1\n> 2\n> 3\n> "
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}

	gotErr := errOutput.String()
	if gotErr != "" {
		t.Fatalf("expected no err, got %q", err)
	}
}

func TestRunFile_ReadsFileFromPathAndOutputsIt(t *testing.T) {
	t.Parallel()
	var output bytes.Buffer
	var errOutput bytes.Buffer

	g, err := golox.NewGolox(
		golox.WithOutput(&output),
		golox.WithStderr(&errOutput),
	)

	if err != nil {
		t.Fatal(err)
	}

	err, exitCode := g.RunFile("testdata/123.txt")
	if err != nil {
		t.Fatal(err)
	}
	if exitCode != 0 {
		t.Fatalf("expected 0 exit code, got %d", exitCode)
	}

	if errOutput.String() != "" {
		t.Fatal("expected no error with valid input")
	}

	got := output.String()
	want := "1\n2\n3\n\n"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}
