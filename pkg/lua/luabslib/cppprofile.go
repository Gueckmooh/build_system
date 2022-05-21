package luabslib

//go:generate go run ./gen -i ./cppprofile.go -c CPPProfile -T ./gen/templates -P luabslib -o cppprofile_gen.go

import (
	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

type CPPProfile struct {
	FDialect      string
	FBuildOptions []string
	FLinkOptions  []string
}

func (p *CPPProfile) Dialect(d string) {
	p.FDialect = d
}

func (p *CPPProfile) AddBuildOptions(bo ...string) {
	p.FBuildOptions = append(p.FBuildOptions, bo...)
}

func (p *CPPProfile) AddLinkOptions(bo ...string) {
	p.FLinkOptions = append(p.FLinkOptions, bo...)
}

func NewCPPProfileLoader(ret **CPPProfile) lua.LGFunction {
	return __NewCPPProfileLoader(ret)
}

func RegisterCPPProfileType(L *lua.LState) {
	__RegisterCPPProfileType(L)
}

func NewCPPProfile() *CPPProfile {
	return &CPPProfile{
		FDialect:      "",
		FBuildOptions: []string{},
		FLinkOptions:  []string{},
	}
}

func ConvertLuaCPPProfileToCPPProfile(cpp *CPPProfile) *project.CPPProfile {
	ccpp := project.NewCPPProfile()
	ccpp.SetDialectFromString(cpp.FDialect)
	ccpp.BuildOptions = cpp.FBuildOptions
	ccpp.LinkOptions = cpp.FLinkOptions
	return ccpp
}
