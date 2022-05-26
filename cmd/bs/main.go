package main

import (
	"fmt"
	"os"

	"github.com/gueckmooh/bs/pkg/argparse"
	log "github.com/gueckmooh/bs/pkg/logging"
	"github.com/gueckmooh/bs/pkg/version"
)

const ProgramName = "bs"

type Options struct {
	parser *argparse.Parser

	debug   *bool
	verbose *bool
	version *bool

	buildOptions BuildOptions
	cleanOptions CleanOptions
}

func (opts *Options) init() {
	opts.parser = argparse.NewParser("bs", "Manages the build system")
	opts.debug = opts.parser.Flag("", "debug", &argparse.Options{
		Required: false,
		Help:     "Print debugging information in addition to normal processing.",
	})
	opts.verbose = opts.parser.Flag("", "verbose", &argparse.Options{
		Required: false,
		Help:     "Make more noise.",
	})
	opts.version = opts.parser.Flag("", "version", &argparse.Options{
		Required: false,
		Help:     "Prints the version",
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

	if *opts.version {
		v, err := version.GetVersion()
		if err != nil {
			return err
		}
		if v.CommitsAhead > 0 {
			fmt.Printf("%s: version %s commit %s (%d commits ahead)\n", ProgramName, v, v.Commit, v.CommitsAhead)
		} else {
			fmt.Printf("%s: version %s\n", ProgramName, v)
		}
		fmt.Printf("Build on %s\n", v.BuildTime)
		return nil
	}

	if *opts.debug {
		log.SetDebugLogging(true)
		log.SetVerboseLogging(true)
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
