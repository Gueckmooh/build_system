package main

import (
	"fmt"
	"os"

	"github.com/gueckmooh/bs/pkg/build"
	projectutils "github.com/gueckmooh/bs/pkg/project_utils"
)

func tryMain() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	proj, err := projectutils.GetProject(cwd)
	if err != nil {
		return err
	}

	builder := build.NewBuilder(proj, build.BuildLib)
	err = builder.BuildComponent("hello_lib")
	if err != nil {
		return err
	}

	return nil
}

func String(s string) *string {
	return &s
}

func main() {
	if err := tryMain(); err != nil {
		fmt.Printf("Fatal error:\n%s\n", err.Error())
	}
}
