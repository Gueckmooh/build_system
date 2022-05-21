package luabslib

import (
	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

//go:generate go run ./gen -i ./profile.go -c Profile -T ./gen/templates -P luabslib -o profile_gen.go

type Profile struct {
	FName    string
	FSources []string
	FCPP     *CPPProfile
}

func NewProfile(name string) *Profile {
	return &Profile{
		FName: name,
		FCPP:  NewCPPProfile(),
	}
}

func (p *Profile) CPP() *CPPProfile {
	return p.FCPP
}

func NewProfileLoader(ret **Profile) lua.LGFunction {
	return __NewProfileLoader(ret)
}

func RegisterProfileType(L *lua.LState) {
	__RegisterProfileType(L)
}

func ConvertLuaProfileToProfile(prof *Profile) *project.Profile {
	pprof := project.NewProfile(prof.FName)
	pprof.SetCPPProfile(ConvertLuaCPPProfileToCPPProfile(prof.FCPP))
	pprof.Sources = functional.ListMap(prof.FSources,
		func(s string) project.FilesPattern { return project.FilesPattern(s) })
	return pprof
}
