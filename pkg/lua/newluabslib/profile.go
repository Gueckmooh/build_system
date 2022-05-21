package newluabslib

import lua "github.com/yuin/gopher-lua"

//go:generate go run ./gen -i ./profile.go -c Profile -T ./gen/templates -P newluabslib -o profile_gen.go

type Profile struct {
	name    string
	sources []string
	_CPP    *CPPProfile
}

func NewProfile(name string) *Profile {
	return &Profile{
		_CPP: NewCPPProfile(name),
	}
}

func (p *Profile) CPP() *CPPProfile {
	return p._CPP
}

func NewProfileLoader(ret **Profile) lua.LGFunction {
	return __NewProfileLoader(ret)
}

func RegisterProfileType(L *lua.LState) {
	__RegisterProfileType(L)
}
