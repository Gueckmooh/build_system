package main

import (
	"fmt"
	"os"

	"github.com/gueckmooh/bs/pkg/lua"
)

func tryMain() error {
	C := lua.NewLuaContext()
	defer C.Close()
	proj, err := C.ReadProjectFile("bs_project.lua")
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", proj)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	files, err := proj.GetComponentFiles(cwd)
	if err != nil {
		return err
	}

	fmt.Printf("files: %#v\n", files)

	components, err := C.ReadComponentFiles(files)
	if err != nil {
		return err
	}

	for _, c := range components {
		fmt.Printf("component: %#v\n", c)
	}

	return nil
}

func main() {
	if err := tryMain(); err != nil {
		fmt.Printf("Fatal error:\n%s\n", err.Error())
	}
}
