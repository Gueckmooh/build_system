package gcc

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/gueckmooh/bs/pkg/common/colors"
	"github.com/gueckmooh/bs/pkg/functional"
	log "github.com/gueckmooh/bs/pkg/logging"
	"github.com/gueckmooh/bs/pkg/project"
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
	gpp          bool
	debugLevel   DebugLevel
	includes     []string
	libDirs      []string
	libs         []string
	targetKind   TargetKind
	dialect      int8
	buildOptions []string
	linkOptions  []string
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

func WithDialect(dialect int8) GCCOption {
	return func(g *GCC) {
		g.dialect = dialect
	}
}

func WithBuildOption(s string) GCCOption {
	return func(g *GCC) {
		g.buildOptions = append(g.buildOptions, s)
	}
}

func WithLinkOption(s string) GCCOption {
	return func(g *GCC) {
		g.linkOptions = append(g.linkOptions, s)
	}
}

func NewGPP(opts ...GCCOption) *GCC {
	gcc := &GCC{
		gpp:        true,
		debugLevel: DebugLevelO0,
		includes:   []string{},
		targetKind: targetExe,
		dialect:    project.DialectCPPUnknown,
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
	log.Log.Printf("%s\n", shellescape.QuoteCommand(cmd))
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

	dialectopt := gcc.getDialectOption()
	if dialectopt != "" {
		cmd = append(cmd, dialectopt)
	}

	for _, v := range gcc.buildOptions {
		cmd = append(cmd, v)
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

	fmt.Printf("Compiling %s%s%s\n", colors.StyleBold, source, colors.StyleReset)
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

	for _, v := range gcc.linkOptions {
		cmd = append(cmd, v)
	}

	fmt.Printf("Linking %s%s%s\n", colors.StyleBold, target, colors.StyleReset)
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

func (g *GCC) getDialectOption() string {
	if g.gpp {
		switch g.dialect {
		case project.DialectCPP98:
			return "-std=c++98"
		case project.DialectCPP03:
			return "-std=c++03"
		case project.DialectCPP11:
			return "-std=c++11"
		case project.DialectCPP0x:
			return "-std=c++0x"
		case project.DialectCPP14:
			return "-std=c++14"
		case project.DialectCPP1y:
			return "-std=c++1y"
		case project.DialectCPP17:
			return "-std=c++17"
		case project.DialectCPP1z:
			return "-std=c++1z"
		case project.DialectCPP20:
			return "-std=c++20"
		case project.DialectCPP2a:
			return "-std=c++2a"
		case project.DialectCPP23:
			return "-std=c++23"
		case project.DialectCPP2b:
			return "-std=c++2b"
		case project.DialectCPPGNU98:
			return "-std=gnu++98"
		case project.DialectCPPGNU03:
			return "-std=gnu++03"
		case project.DialectCPPGNU11:
			return "-std=gnu++11"
		case project.DialectCPPGNU0x:
			return "-std=gnu++0x"
		case project.DialectCPPGNU14:
			return "-std=gnu++14"
		case project.DialectCPPGNU1y:
			return "-std=gnu++1y"
		case project.DialectCPPGNU17:
			return "-std=gnu++17"
		case project.DialectCPPGNU1z:
			return "-std=gnu++1z"
		case project.DialectCPPGNU20:
			return "-std=gnu++20"
		case project.DialectCPPGNU2a:
			return "-std=gnu++2a"
		case project.DialectCPPGNU23:
			return "-std=gnu++23"
		case project.DialectCPPGNU2b:
			return "-std=gnu++2b"
		}
	}
	return ""
}
