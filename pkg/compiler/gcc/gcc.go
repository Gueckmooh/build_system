package gcc

import (
	"bytes"
	"fmt"
	"os"
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

type TargetKind int8

const (
	targetExe TargetKind = iota
	targetLib
)

type GCC struct {
	gpp        bool
	debugLevel DebugLevel
	includes   []string
	libDirs    []string
	libs       []string
	targetKind TargetKind
}

type GCCOption func(*GCC)

func TargetLib(g *GCC) {
	g.targetKind = targetLib
}

func WithInclude(include string) GCCOption {
	return func(g *GCC) {
		g.includes = append(g.includes, include)
	}
}

func WithLibDir(libDir string) GCCOption {
	return func(g *GCC) {
		g.libDirs = append(g.libDirs, libDir)
	}
}

func WithLib(lib string) GCCOption {
	return func(g *GCC) {
		g.libs = append(g.libs, lib)
	}
}

func NewGPP(opts ...GCCOption) *GCC {
	gcc := &GCC{
		gpp:        true,
		debugLevel: DebugLevelO0,
		includes:   []string{},
		targetKind: targetExe,
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

	if gcc.targetKind == targetLib {
		cmd = append(cmd, "-fPIC")
	}

	cmd = append(cmd, includesOpts...)

	cmd = append(cmd, "-c")

	cmd = append(cmd, source)

	cmd = append(cmd, []string{"-o", target}...)

	fmt.Printf("Compiling %s\n", source)
	_, errs, err := runCommand(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", errs)
		return fmt.Errorf("Error while compiling file %s\n\t%s", source, err.Error())
	}
	return nil
}

func (gcc *GCC) LinkFiles(target string, sources ...string) error {
	var cmd []string
	if gcc.gpp {
		cmd = append(cmd, GPPExec)
	} else {
		cmd = append(cmd, GCCExec)
	}

	if gcc.targetKind == targetLib {
		cmd = append(cmd, "-shared")
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
	_, errs, err := runCommand(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", errs)
		return fmt.Errorf("Error while linking file %s\n\t%s", target, err.Error())
	}
	return nil
}

func (gcc *GCC) GetFileDependencies(target, source string) (string, []string, error) {
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

	outs, errs, err := runCommand(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", errs)
		return "", nil, fmt.Errorf("Error while compiling file %s\n\t%s", source, err.Error())
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
