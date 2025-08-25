package golox

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/taylorlowery/lox/internal/ast"
	"github.com/taylorlowery/lox/internal/parser"
	"github.com/taylorlowery/lox/internal/scanner"
	"github.com/taylorlowery/lox/internal/token"
)

type Golox struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	hadErr bool
}

type option func(*Golox) error

// NewGolox returns a new instance of NewGolox
// with the input configured to Stdin
// and output to Stdout
func NewGolox(opts ...option) (*Golox, error) {
	g := &Golox{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
	for _, opt := range opts {
		err := opt(g)
		if err != nil {
			return nil, err
		}
	}
	return g, nil
}

// WithOutput configures a Golox instance to use a given writer as output
func WithOutput(w io.Writer) option {
	return func(g *Golox) error {
		if w == nil {
			return errors.New("nil output writer")
		}
		g.stdout = w
		return nil
	}
}

// WithInput configures a Golox instance to use a given reader as input
func WithInput(r io.Reader) option {
	return func(g *Golox) error {
		if r == nil {
			return errors.New("nil input reader")
		}
		g.stdin = r
		return nil
	}
}

// WithStderr configures a given Golox instance to use a given writer for error output
func WithStderr(w io.Writer) option {
	return func(g *Golox) error {
		if w == nil {
			return errors.New("nil output writer")
		}
		g.stderr = w
		return nil
	}
}

func (g Golox) run(source string) {
	scanner := scanner.NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		g.Error(err.Line, err.Message)
		return
	}

	p := parser.NewParser(tokens)
	expr, parseErr := p.Parse()
	if parseErr != nil {
		g.Error(0, parseErr.Error())
	}

	astPrinter, printerError := ast.NewAstPrinter()
	if printerError != nil {
		panic(printerError)
	}

	fmt.Println(astPrinter.PrintAst(expr))
}

// RunFile reads a file at a given path,
// parses it and executes it as Lox
func (g Golox) RunFile(filepath string) (error, int) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return err, 65
	}
	g.run(string(bytes))
	return nil, 0
}

// RunPrompt starts a loop over the Golox's input,
// parsing its input line by line and executing it as Lox
func (g Golox) RunPrompt() error {
	// thusly named to differentiate from my own scanner
	bufioScanner := bufio.NewScanner(g.stdin)
	for {
		fmt.Fprint(g.stdout, "> ")
		if !bufioScanner.Scan() {
			break
		}
		line := bufioScanner.Text()
		if line == "" {
			continue
		}
		g.run(line)
		g.hadErr = false
	}
	return bufioScanner.Err()
}

func (g Golox) Error(line int, message string) {
	g.report(line, "", message)
}

func (g Golox) TokenError(t token.Token, message string) {
	if t.TokenType == token.EOF {
		g.report(t.Line, " at end", message)
	} else {
		g.report(t.Line, " at "+t.Lexeme, message)
	}
}

func (g Golox) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line: %d] Error %s: %s\n", line, where, message)
	g.hadErr = true
}

func (g Golox) HadError() bool {
	return g.hadErr
}
