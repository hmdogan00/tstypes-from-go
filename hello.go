package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("gotype", "Converts go source code into typescript types")
	input := parser.String("i", "input", &argparse.Options{Required: true, Help: "Input golang directory/file to type"})
	output := parser.String("o", "output", &argparse.Options{Required: true, Help: "Output file name to export the types"})

	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	fmt.Println(*input, *output)
}
