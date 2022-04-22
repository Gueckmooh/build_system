package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gueckmooh/bs/pkg/argparse"
)

type Method struct {
	Name   string `xml:"name,attr"`
	Kind   string `xml:"kind,attr"`
	Type   string `xml:"type,attr"`
	Target string `xml:"target,attr"`
}

type Field struct {
	Name    string `xml:"name,attr"`
	Type    string `xml:"type,attr"`
	Private bool   `xml:"private,attr"`
}

type Table struct {
	XMLName xml.Name `xml:"table"`
	Name    string   `xml:"name,attr"`
	Methods struct {
		Method []Method `xml:"method"`
	} `xml:"methods"`
	Fields struct {
		Field []Field `xml:"field"`
	} `xml:"fields"`
}

func tryMain() error {
	parser := argparse.NewParser("gen", "Generate lua binding")
	inputFile := parser.String("i", "input-file", &argparse.Options{
		Required: true,
		Help:     "The xml input file",
	})
	outputFile := parser.String("o", "output-file", &argparse.Options{
		Required: false,
		Help:     "The go output file",
	})
	packageName := parser.String("", "package", &argparse.Options{
		Required: true,
		Help:     "The package to use",
	})
	publicInterfaceName := parser.String("", "public-interface", &argparse.Options{
		Help: "The name of the public interface",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return err
	}
	var t Table
	err = xml.Unmarshal(data, &t)
	if err != nil {
		return err
	}
	tg := NewTableGenerator(&t, WithPackageName(*packageName), WithPublicInterface(*publicInterfaceName))
	if len(*outputFile) > 0 {
		err := ioutil.WriteFile(*outputFile, []byte(tg.GenFile()), 0o600)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(tg.GenFile())
	}
	return nil
}

func main() {
	err := tryMain()
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
