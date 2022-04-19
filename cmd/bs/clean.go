package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gueckmooh/bs/pkg/argparse"
	"github.com/gueckmooh/bs/pkg/fsutil"
	"github.com/gueckmooh/bs/pkg/lua"
	"github.com/gueckmooh/bs/pkg/project"
)

type CleanOptions struct {
	command *argparse.Command

	name *argparse.PosStringResult
}

func (opts *CleanOptions) init(parser *argparse.Parser) {
	opts.command = parser.NewCommand("clean", "Clean project or component")
}

func (opts *CleanOptions) happened() bool {
	return opts.command.Happened()
}

func tryCleanMain(opts Options) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	projFile, err := fsutil.FindFileUpstream(project.ProjectConfigFile, cwd)
	if err != nil {
		return err
	}

	cwd = filepath.Dir(projFile)
	err = os.Chdir(cwd)
	if err != nil {
		return err
	}

	C := lua.NewLuaContext()
	defer C.Close()
	proj, err := C.GetProject(cwd)
	if err != nil {
		return err
	}

	fmt.Printf("Removing dir %s\n", proj.Config.GetBuildDirectory(true))
	err = os.RemoveAll(proj.Config.GetBuildDirectory(true))
	if err != nil {
		return err
	}

	return nil
}

func cleanMain(opts Options) error {
	err := tryCleanMain(opts)
	if err != nil {
		return fmt.Errorf("Error while building components:\n  %s", err.Error())
	}

	return nil
}
