package main

import (
	"fmt"
	"os"

	"github.com/gueckmooh/argparse"
	alist "github.com/gueckmooh/bs/pkg/adjacency_list"
	"github.com/gueckmooh/bs/pkg/build"
	"github.com/gueckmooh/bs/pkg/project"
	projectutils "github.com/gueckmooh/bs/pkg/project_utils"
)

type BuildOptions struct {
	command *argparse.Command

	name          *argparse.PosStringResult
	buildUpstream *bool
}

func (opts *BuildOptions) init(parser *argparse.Parser) {
	opts.command = parser.NewCommand("build", "Build project or component")

	opts.name = opts.command.PosString("component", &argparse.Options{
		Required: false,
		Help:     "The name of the component to build",
	})
	opts.buildUpstream = opts.command.Flag("", "build-upstream", &argparse.Options{
		Required: false,
		Help:     "Instruct to build all the upstream components",
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

	err = proj.ComputeComponentDependencies()
	if err != nil {
		return err
	}
	// os.Exit(0)

	// @todo make sure there is no cycle

	ctb := proj.DefaultTarget
	if ctb == "" {
		if len(*opts.buildOptions.name) == 0 {
			return fmt.Errorf("No component name given")
		}
		ctb = (*opts.buildOptions.name)[0]
	}

	if *opts.buildOptions.buildUpstream {
		err = BuildUpstream(proj, ctb)
		if err != nil {
			return err
		}
	} else {
		builder := build.NewBuilder(proj, ctb)
		err = builder.BuildComponent()
		if err != nil {
			return err
		}
	}

	return nil
}

func BuildUpstream(proj *project.Project, ctb string) error {
	var processNode func(alist.VertexDescriptor) error
	g := proj.ComponentDeps.G
	processNode = func(v alist.VertexDescriptor) error {
		oe, err := g.OutEdges(v)
		if err != nil {
			return err
		}

		for _, ed := range oe {
			target, err := g.Target(ed)
			if err != nil {
				return err
			}
			err = processNode(target)
			if err != nil {
				return err
			}
		}
		builder := build.NewBuilder(proj, g.GetVertexAttribute(v).Name)
		err = builder.BuildComponent()
		if err != nil {
			return err
		}
		return nil
	}
	c, err := proj.GetComponent(ctb)
	if err != nil {
		return err
	}
	return processNode(proj.ComponentDeps.Vmap[c])
}

func buildMain(opts Options) error {
	err := tryBuildMain(opts)
	if err != nil {
		return fmt.Errorf("Error while building components:\n  %s", err.Error())
	}

	return nil
}
