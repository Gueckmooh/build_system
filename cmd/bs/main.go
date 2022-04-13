package main

import (
	"fmt"
	"os"

	"github.com/gueckmooh/bs/pkg/argparse"
	log "github.com/gueckmooh/bs/pkg/logging"
)

type Options struct {
	parser *argparse.Parser

	debug   *bool
	verbose *bool
	// directory *string

	buildOptions BuildOptions
	cleanOptions CleanOptions
}

func (opts *Options) init() {
	opts.parser = argparse.NewParser("bs", "Manages the build system")
	opts.debug = opts.parser.Flag("", "debug", &argparse.Options{
		Required: false,
		Help:     "Print debugging information in addition to normal processing.",
		// Default:  false,
	})
	opts.verbose = opts.parser.Flag("", "verbose", &argparse.Options{
		Required: false,
		Help:     "Make more noise.",
		// Default:  false,
	})
	opts.buildOptions.init(opts.parser)
	opts.cleanOptions.init(opts.parser)
}

func tryMain() error {
	var opts Options
	opts.init()

	err := opts.parser.Parse(os.Args)
	if err != nil {
		return fmt.Errorf("Fails to parse options:\n  %s\n", err.Error())
	}

	if *opts.debug {
		log.SetDebugLogging(true)
	}
	if *opts.verbose {
		log.SetVerboseLogging(true)
	}

	if opts.buildOptions.happened() {
		return buildMain(opts)
	} else if opts.cleanOptions.happened() {
		return cleanMain(opts)
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
