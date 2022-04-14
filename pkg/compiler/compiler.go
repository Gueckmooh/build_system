package compiler

import (
	"github.com/gueckmooh/bs/pkg/compiler/gcc"
	"github.com/gueckmooh/bs/pkg/project"
)

type Compiler interface {
	CompileFile(target, source string) error
	LinkFiles(target string, sources ...string) error
	GetFileDependencies(target, source string) (string, []string, error)
}

const (
	targetExe int8 = iota
	targetLib
)

type compilerOption struct {
	includeDirectories []string
	libraryDirectories []string
	libraries          []string
	forCPP             bool
	targetKind         int8
	cppDialect         int8
	buildOptions       []string
}

type CompilerOption func(*compilerOption)

func WithIncludeDirectory(incDir string) CompilerOption {
	return func(co *compilerOption) {
		co.includeDirectories = append(co.includeDirectories, incDir)
	}
}

func WithLibraryDirectory(libDir string) CompilerOption {
	return func(co *compilerOption) {
		co.libraryDirectories = append(co.libraryDirectories, libDir)
	}
}

func WithLibrary(lib string) CompilerOption {
	return func(co *compilerOption) {
		co.libraries = append(co.libraries, lib)
	}
}

func ForCPP(co *compilerOption) {
	co.forCPP = true
}

func TargetLib(co *compilerOption) {
	co.targetKind = targetLib
}

func TargetExe(co *compilerOption) {
	co.targetKind = targetExe
}

func WithCPPDIalect(dialect int8) CompilerOption {
	return func(co *compilerOption) {
		co.cppDialect = dialect
	}
}

func WithBuildOption(s string) CompilerOption {
	return func(co *compilerOption) {
		co.buildOptions = append(co.buildOptions, s)
	}
}

func NewCompiler(opts ...CompilerOption) Compiler {
	options := &compilerOption{
		forCPP:     false,
		cppDialect: project.DialectCPPUnknown,
	}
	for _, opt := range opts {
		opt(options)
	}
	return options.newGCCCompiler()
}

func (co *compilerOption) newGCCCompiler() Compiler {
	var opts []gcc.GCCOption
	for _, v := range co.includeDirectories {
		opts = append(opts, gcc.WithInclude(v))
	}
	for _, v := range co.libraryDirectories {
		opts = append(opts, gcc.WithLibDir(v))
	}
	for _, v := range co.libraries {
		opts = append(opts, gcc.WithLib(v))
	}
	switch co.targetKind {
	case targetLib:
		opts = append(opts, gcc.TargetLib)
	}
	if co.cppDialect != project.DialectCPPUnknown {
		opts = append(opts, gcc.WithDialect(co.cppDialect))
	}
	for _, v := range co.buildOptions {
		opts = append(opts, gcc.WithBuildOption(v))
	}
	return gcc.NewGPP(opts...)
}
