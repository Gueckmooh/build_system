package gcc

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gueckmooh/bs/pkg/functional"
)

type DebugLevel int8

const (
	DebugLevelO0 DebugLevel = iota
	DebugLevelO1
	DebugLevelO2
	DebugLevelO3
	DebugLevelOg
)

const (
	GPPExec = "g++"
	GCCExec = "gcc"
)

type GCC struct {
	gpp        bool
	debugLevel DebugLevel
	includes   []string
	libDirs    []string
	libs       []string
}

type GCCOption func(*GCC)

func NewGPP(opts ...GCCOption) *GCC {
	gcc := &GCC{
		gpp:        true,
		debugLevel: DebugLevelO0,
		includes:   []string{},
	}
	for _, opt := range opts {
		opt(gcc)
	}

	return gcc
}

func runCommand(cmd []string) error {
	fmt.Println(strings.Join(cmd, " "))
	exe := exec.Command(cmd[0], cmd[1:]...)
	err := exe.Run()
	if err != nil {
		return err
	}
	return nil
}

func (gcc *GCC) CompileFile(target, source string) error {
	var cmd []string
	if gcc.gpp {
		cmd = append(cmd, GPPExec)
	} else {
		cmd = append(cmd, GCCExec)
	}

	includesOpts := functional.ListMap(gcc.includes,
		func(s string) string {
			return "-I" + s
		})

	cmd = append(cmd, includesOpts...)

	cmd = append(cmd, "-c")

	cmd = append(cmd, source)

	cmd = append(cmd, []string{"-o", target}...)

	runCommand(cmd)

	return nil
}

func (gcc *GCC) LinkFile(target string, sources ...string) error {
	var cmd []string
	if gcc.gpp {
		cmd = append(cmd, GPPExec)
	} else {
		cmd = append(cmd, GCCExec)
	}

	libDirOpts := functional.ListMap(gcc.libDirs,
		func(s string) string {
			return "-L" + s
		})
	libOpts := functional.ListMap(gcc.libs,
		func(s string) string {
			return "-l" + s
		})

	cmd = append(cmd, libDirOpts...)
	cmd = append(cmd, libOpts...)

	cmd = append(cmd, sources...)

	cmd = append(cmd, []string{"-o", target}...)

	runCommand(cmd)

	return nil
}
