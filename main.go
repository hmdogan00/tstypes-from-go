package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/akamensky/argparse"
)

func fileToTs(file *ast.File) string {
	result := ""
	ast.Inspect(file, func(n ast.Node) bool {
		genDecl, genOk := n.(*ast.GenDecl)
		funcDecl, funcOk := n.(*ast.FuncDecl)
		if genOk {
			if genDecl.Tok == token.VAR {
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						result += varHandler(*valueSpec)
					}
				}
			}
			if genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							result += interfaceHandler(*typeSpec, *structType)
						}
					}
				}
			}
			if genDecl.Tok == token.CONST {
				if len(genDecl.Specs) > 1 {
					result += enumHandler(genDecl.Specs)
				} else {
					if valueSpec, ok := genDecl.Specs[0].(*ast.ValueSpec); ok {
						result += varHandler(*valueSpec)
					} else {
						fmt.Println("Unhandled const type detected")
						panic(genDecl.Specs[0])
					}
				}
			}
		}
		if funcOk {
			result += functionHandler(*funcDecl)
		}
		return true
	})
	return result
}

func main() {
	argParser := argparse.NewParser("gotstype", "Converts go source code into typescript types")
	input := argParser.String("i", "input", &argparse.Options{Required: true, Help: "Input golang directory/file to type"})
	outputFile := argParser.String("o", "output", &argparse.Options{Required: true, Help: "Output directory name to export the types"})

	argErr := argParser.Parse(os.Args)

	if argErr != nil {
		fmt.Print(argParser.Usage(argErr))
	}

	// if output does not have any extension, add .d.ts to it.
	if !strings.Contains(*outputFile, ".") {
		*outputFile = *outputFile + ".d.ts"
	}

	fset := token.NewFileSet()
	ast, err := parser.ParseDir(fset, *input, nil, 0)

	if err != nil {
		panic(err)
	}
	output := ""

	for _, pkg := range ast {
		for _, file := range pkg.Files {
			output += fileToTs(file)
		}
	}

	os.WriteFile(*outputFile, []byte(output), 0644)
}
