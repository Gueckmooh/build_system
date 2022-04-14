package main

import (
	"fmt"
	"os"
	"path/filepath"

	alist "github.com/gueckmooh/bs/pkg/adjacency_list"
	"github.com/gueckmooh/bs/pkg/argparse"
	"github.com/gueckmooh/bs/pkg/build"
	"github.com/gueckmooh/bs/pkg/common/colors"
	"github.com/gueckmooh/bs/pkg/fsutil"
	log "github.com/gueckmooh/bs/pkg/logging"
	"github.com/gueckmooh/bs/pkg/project"
	projectutils "github.com/gueckmooh/bs/pkg/project_utils"
)

type BuildOptions struct {
	command *argparse.Command

	name          *argparse.PosStringResult
	buildUpstream *bool
	directory     *string
	alwaysBuild   *bool
	profile       *string
	platform      *string
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
	opts.directory = opts.command.String("C", "directory", &argparse.Options{
		Validate: func(args []string) error {
			for _, s := range args {
				if s == "" {
					continue
				}
				stats, err := os.Stat(s)
				if err != nil {
					return err
				}
				if !stats.IsDir() {
					return fmt.Errorf("%s is not a directory", s)
				}
			}
			return nil
		},
		Help: "Change to directory dir before reading the bsfiles or doing anything else.",
	})
	opts.alwaysBuild = opts.command.Flag("B", "always-build", &argparse.Options{
		Required: false,
		Help:     "Unconditionally build all targets.",
	})
	opts.profile = opts.command.String("p", "profile", &argparse.Options{
		Required: false,
		Help:     "Use selected profile for build.",
	})
	opts.platform = opts.command.String("P", "platform", &argparse.Options{
		Required: false,
		Help:     "Use selected platform for build.",
	})
}

func (opts *BuildOptions) happened() bool {
	return opts.command.Happened()
}

func tryBuildMainInDirectory(directory string, opts Options) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Log.Printf("Movind to directory '%s'\n", directory)
	err = os.Chdir(directory)
	if err != nil {
		return err
	}
	err = tryBuildMain(opts)
	if err != nil {
		return err
	}
	log.Log.Printf("Exiting directory '%s'\n", directory)
	err = os.Chdir(cwd)
	if err != nil {
		return err
	}
	return nil
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

	var bops []build.BuildOption
	if *opts.buildOptions.alwaysBuild {
		bops = append(bops, build.WithAlwaysBuild)
	}
	profilestr := "unspecified"
	platformstr := "unspecified"
	if *opts.buildOptions.profile != "" {
		bops = append(bops, build.WithProfile(*opts.buildOptions.profile))
		profilestr = *opts.buildOptions.profile
	} else if proj.DefaultProfile != "" {
		bops = append(bops, build.WithProfile(proj.DefaultProfile))
		profilestr = proj.DefaultProfile
	}
	if *opts.buildOptions.platform != "" {
		bops = append(bops, build.WithPlatform(*opts.buildOptions.platform))
		platformstr = *opts.buildOptions.platform
	} else if proj.DefaultPlatform != "" {
		bops = append(bops, build.WithPlatform(proj.DefaultPlatform))
		platformstr = proj.DefaultPlatform
	}

	log.Log.Printf("%sInfo:%s build configured for %s profile, %s platform...\n",
		colors.ColorCyan, colors.ColorReset, profilestr, platformstr)

	if *opts.buildOptions.buildUpstream {
		err = BuildUpstream(proj, ctbs[0], bops)
		if err != nil {
			return err
		}
	} else {
		for _, ctb := range ctbs {
			builder, err := build.NewBuilder(proj, ctb, bops...)
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

func BuildUpstream(proj *project.Project, ctb string, bops []build.BuildOption) error {
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
		builder, err := build.NewBuilder(proj, g.GetVertexAttribute(v).Name, bops...)
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
	var err error
	if len(*opts.buildOptions.directory) > 0 {
		err = tryBuildMainInDirectory(*opts.buildOptions.directory, opts)
	} else {
		err = tryBuildMain(opts)
	}
	if err != nil {
		return fmt.Errorf("Error while building components:\n  %s", err.Error())
	}

	return nil
}
