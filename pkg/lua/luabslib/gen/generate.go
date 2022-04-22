package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/gueckmooh/bs/pkg/argparse"
)

type Method struct {
	Name   string `xml:"name,attr"`
	Kind   string `xml:"kind,attr"`
	Type   string `xml:"type,attr"`
	Target string `xml:"target,attr"`
}

type Field struct {
	Name             string `xml:"name,attr"`
	Type             string `xml:"type,attr"`
	Private          bool   `xml:"private,attr"`
	DefaultFromParam int    `xml:"default_from_param,attr"`
}

type Constructor struct {
	Type string `xml:"type,attr"`
}

type Table struct {
	XMLName xml.Name `xml:"table"`
	Require []struct {
		Name string `xml:"name,attr"`
	} `xml:"require"`
	Constructor *Constructor `xml:"constructor"`
	Name        string       `xml:"name,attr"`
	Methods     []Method     `xml:"methods>method"`
	Fields      []Field      `xml:"fields>field"`
}

func readXMLFile(filename string) (*Table, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var t Table
	err = xml.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
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

	t, err := readXMLFile(*inputFile)

	for _, req := range t.Require {
		path := path.Join(path.Dir(*inputFile), req.Name)
		tt, err := readXMLFile(path)
		if err != nil {
			return err
		}
		treq := NewTableGenerator(tt)
		Dependencies = append(Dependencies, treq)
	}

	tg := NewTableGenerator(t, WithPackageName(*packageName), WithPublicInterface(*publicInterfaceName))

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
