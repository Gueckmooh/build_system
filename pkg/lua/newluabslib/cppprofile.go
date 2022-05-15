package newluabslib

import lua "github.com/yuin/gopher-lua"

type CPPProfile struct {
	DialectF     string
	BuildOptions []string
	LinkOptions  []string
}

func (p *CPPProfile) Dialect(d string) {
	p.DialectF = d
}

func (p *CPPProfile) AddBuildOptions(bo ...string) {
	p.BuildOptions = append(p.BuildOptions, bo...)
}

func (p *CPPProfile) AddLinkOptions(bo ...string) {
	p.LinkOptions = append(p.LinkOptions, bo...)
}

func NewCPPProfileLoader(ret **CPPProfile) lua.LGFunction {
	return __NewCPPProfileLoader(ret)
}

func NewCPPProfile(dialect string) *CPPProfile {
	return &CPPProfile{DialectF: dialect}
}
