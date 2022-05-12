package main

import (
	"fmt"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"

	"github.com/gueckmooh/bs/pkg/argparse"
)

var packageName string

func tryMain() error {
	optparser := argparse.NewParser("gen", "Generate lua bindings")
	inputFile := optparser.String("i", "input-file", &argparse.Options{
		Required: true,
		Help:     "The go input file",
	})
	className := optparser.String("c", "class-name", &argparse.Options{
		Required: true,
		Help:     "The class for which we generate bindings",
	})
	tmplDir := optparser.String("T", "templates-dir", &argparse.Options{
		Required: true,
		Help:     "The directory containing the templates to parse",
	})
	outputFile := optparser.String("o", "output-file", &argparse.Options{
		Required: false,
		Help:     "The file in which to write the output",
	})
	pckgName := optparser.String("P", "package-name", &argparse.Options{
		Required: true,
		Help:     "The name of the package",
	})

	err := optparser.Parse(os.Args)
	if err != nil {
		return err
	}

	templateDir = *tmplDir
	packageName = *pckgName

	fset := token.NewFileSet()
	data, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fset, *inputFile, data, parser.AllErrors)
	if err != nil {
		return err
	}

	printer.Fprint(os.Stdout, fset, file)

	classes := getClasses(file)
	var class *Class
	for _, c := range classes {
		if c.Name == *className {
			class = c
		}
	}
	if class == nil {
		return fmt.Errorf("Could not find class %s", *className)
	}

	getMethodsForClass(class, file)

	SetFunctionNameBundle(class)

	fmt.Printf("%s\n\n", class)

	body := MustExecuteTemplate("main.gotmpl", class)
	ndata, err := format.Source([]byte(body))
	if err != nil {
		fmt.Println(body)
		return err
	}

	if len(*outputFile) == 0 {
		fmt.Println(string(ndata))
	} else {
		ioutil.WriteFile(*outputFile, ndata, 0o600)
	}

	return nil
}

func main() {
	if err := tryMain(); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
