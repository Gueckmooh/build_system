package gcc

import (
	"bytes"
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

func runCommand(cmd []string) (string, string, error) {
	exe := exec.Command(cmd[0], cmd[1:]...)
	var outb, errb bytes.Buffer
	exe.Stdout = &outb
	exe.Stderr = &errb
	err := exe.Run()
	return outb.String(), errb.String(), err
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

	fmt.Printf("Compiling %s\n", source)
	_, _, err := runCommand(cmd)
	return err
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

	fmt.Printf("Linking %s\n", target)
	_, _, err := runCommand(cmd)
	return err
}

func (gcc *GCC) GetBuildInfoForFile(target, source string) (string, []string, error) {
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

	cmd = append(cmd, "-MM")

	cmd = append(cmd, source)

	cmd = append(cmd, []string{"-MT", target}...)

	outs, _, err := runCommand(cmd)
	if err != nil {
		return "", nil, err
	}
	return ParseMOutput(outs)
}

func ParseMOutput(o string) (string, []string, error) {
	o = strings.ReplaceAll(o, "\\\n", "")
	os := strings.Split(o, ":")
	if len(os) < 2 {
		return "", nil, fmt.Errorf("Error while parsing M output")
	}
	target := os[0]
	sources := strings.Split(strings.TrimSpace(os[1]), " ")
	sources = functional.ListFilter(sources, func(s string) bool {
		return s != ""
	})
	return target, sources, nil
}
