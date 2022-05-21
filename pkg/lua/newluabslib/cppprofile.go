package newluabslib

//go:generate go run ./gen -i ./cppprofile.go -c CPPProfile -T ./gen/templates -P newluabslib -o cppprofile_gen.go

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

func RegisterCPPProfileType(L *lua.LState) {
	__RegisterCPPProfileType(L)
}

func NewCPPProfile(dialect string) *CPPProfile {
	return &CPPProfile{DialectF: dialect}
}
