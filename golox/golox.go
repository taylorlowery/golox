package golox

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

var HadError bool = false

type Golox struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
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
	// scanner := new Scanner(source)
	// tokens := scanner.ScanTokens()

	// for token := range tokens {
	// 	fmt.Fprint(g.output, token)
	// }
	fmt.Fprintln(g.stdout, source)
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
		HadError = false
	}
	return bufioScanner.Err()
}

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line: %d] Error %s: %s\n", line, where, message)
	HadError = true
}
