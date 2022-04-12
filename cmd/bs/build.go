package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gueckmooh/argparse"
	alist "github.com/gueckmooh/bs/pkg/adjacency_list"
	"github.com/gueckmooh/bs/pkg/build"
	"github.com/gueckmooh/bs/pkg/common/colors"
	"github.com/gueckmooh/bs/pkg/fsutil"
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

	projFile, err := fsutil.FindFileUpstream(project.ProjectConfigFile, cwd)
	if err != nil {
		return err
	}
	oldcwd := cwd
	cwd = filepath.Dir(projFile)
	err = os.Chdir(cwd)
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
	// @todo make sure there is no cycle

	ctbs := *opts.buildOptions.name
	if len(ctbs) == 0 {
		if proj.DefaultTarget != "" {
			ctbs = append(ctbs, proj.DefaultTarget)
		} else {
			compFile, err := fsutil.FindFileUpstream(project.ComponentConfigFile, oldcwd)
			if err == nil {
				ctb, err := proj.GetComponentByPath(filepath.Dir(compFile))
				if err == nil {
					ctbs = append(ctbs, ctb.Name)
				}
			}
		}
	}
	if len(ctbs) == 0 {
		return fmt.Errorf("No component name given")
	}

	if len(ctbs) > 1 {
		if *opts.buildOptions.buildUpstream {
			fmt.Fprintf(os.Stderr, "%sWarning%s: several components given, build upstream ignored\n",
				colors.ColorYellow, colors.ColorReset)
			*opts.buildOptions.buildUpstream = false
		}
	}
	if *opts.buildOptions.buildUpstream {
		err = BuildUpstream(proj, ctbs[0])
		if err != nil {
			return err
		}
	} else {
		for _, ctb := range ctbs {
			builder, err := build.NewBuilder(proj, ctb)
			if err != nil {
				return err
			}
			err = builder.BuildComponent()
			if err != nil {
				return err
			}
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
		builder, err := build.NewBuilder(proj, g.GetVertexAttribute(v).Name)
		if err != nil {
			return err
		}
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
