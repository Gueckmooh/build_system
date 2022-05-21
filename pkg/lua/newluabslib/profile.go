package newluabslib

import lua "github.com/yuin/gopher-lua"

//go:generate go run ./gen -i ./profile.go -c Profile -T ./gen/templates -P newluabslib -o profile_gen.go

type Profile struct {
	FName    string
	FSources []string
	FCPP     *CPPProfile
}

func NewProfile(name string) *Profile {
	return &Profile{
		FCPP: NewCPPProfile(),
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
