package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: generate-ast <output directory>")
		os.Exit(64)
	}

	//outputDir := os.Args[1]

}
