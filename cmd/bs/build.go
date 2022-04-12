package main

import (
	"fmt"
	"os"

	"github.com/gueckmooh/argparse"
	"github.com/gueckmooh/bs/pkg/build"
	projectutils "github.com/gueckmooh/bs/pkg/project_utils"
)

type BuildOptions struct {
	command *argparse.Command

	name *argparse.PosStringResult
}

func (opts *BuildOptions) init(parser *argparse.Parser) {
	opts.command = parser.NewCommand("build", "Build project or component")

	opts.name = opts.command.PosString("component", &argparse.Options{
		Required: false,
		Help:     "The name of the component to build",
	})
}

func (opts *BuildOptions) happened() bool {
	return opts.command.Happened()
}

func tryBuildMain(opts Options) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	proj, err := projectutils.GetProject(cwd)
	if err != nil {
		return err
	}

	ctb := proj.DefaultTarget
	if ctb == "" {
		if len(*opts.buildOptions.name) == 0 {
			return fmt.Errorf("No component name given")
		}
		ctb = (*opts.buildOptions.name)[0]
	}

	builder := build.NewBuilder(proj, ctb)
	err = builder.BuildComponent()
	if err != nil {
		return err
	}

	return nil
}

func buildMain(opts Options) error {
	err := tryBuildMain(opts)
	if err != nil {
		return fmt.Errorf("Error while building components:\n  %s", err.Error())
	}

	return nil
}
