package main

import (
	"fmt"
	"os"

	"github.com/taylorlowery/lox/golox"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("usage: golox [script]")
		os.Exit(64)
	}
	g, err := golox.NewGolox()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(os.Args) == 2 {
		err, exitCode := g.RunFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
		}
		if g.HadError() {
			os.Exit(exitCode)
		}
	} else {
		err := g.RunPrompt()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	// just to be thorough
	os.Exit(0)
}
