package luabslib

//go:generate go run ./gen -i definitions/Profile.xml --package luabslib --public-interface "LuaProfile" -o profile_gen.go

import (
	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

func NewProfile(L *lua.LState, name string) *lua.LTable {
	sname := lua.LString(name)
	return NewLuaProfile(L, &sname)
}

func GetLuaCPPProfile(L *lua.LState, T *lua.LTable) *lua.LTable {
	err := __profile_CheckTableIntegrity(L, T)
	if err != nil {
		panic(err)
	}
	return L.GetField(T, "CPP").(*lua.LTable)
}

func ReadProfileFromLuaTable(L *lua.LState, T *lua.LTable) (*project.Profile, error) {
	lprofile, err := GetLuaProfileFromLuaTable(L, T)
	if err != nil {
		return nil, err
	}

	p := project.NewProfile(lprofile.getName())
	p.SetCPPProfile(LuaCPPProfileToCPPProfile((*LuaCPPProfile)(lprofile.getCPP())))
	p.Sources = functional.ListMap(lprofile.getSources(),
		func(s string) project.FilesPattern { return project.FilesPattern(s) })

	return p, nil
}
