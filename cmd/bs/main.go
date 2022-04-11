package main

import (
	"fmt"
	"os"

	"github.com/gueckmooh/argparse"
)

type Options struct {
	parser *argparse.Parser

	buildOptions BuildOptions
}

func (opts *Options) init() {
	opts.parser = argparse.NewParser("bs", "Manages the build system")
	opts.buildOptions.init(opts.parser)
}

func tryMain() error {
	var opts Options
	opts.init()

	err := opts.parser.Parse(os.Args)
	if err != nil {
		return fmt.Errorf("Fails to parse options:\n  %s\n", err.Error())
	}

	if opts.buildOptions.happened() {
		return buildMain(opts)
	}

	return fmt.Errorf("No command given")
}

func String(s string) *string {
	return &s
}

func main() {
	if err := tryMain(); err != nil {
		fmt.Printf("Fatal error:\n%s\n", err.Error())
		os.Exit(1)
	}
}
