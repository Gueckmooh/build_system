package luabslib

//go:generate go run ./gen -i definitions/CPPProfile.xml --package luabslib --public-interface "LuaCPPProfile" -o cppprofile_gen.go

import (
	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

func NewCPPProfile(L *lua.LState) *lua.LTable {
	return NewLuaCPPProfile(L)
}

func ReadCPPProfileFromLuaTable(L *lua.LState, T *lua.LTable) (*project.CPPProfile, error) {
	v, err := GetLuaCPPProfileFromLuaTable(L, T)
	if err != nil {
		return nil, err
	}
	cpp := &project.CPPProfile{
		BuildOptions: v.getBuildOptions(),
		LinkOptions:  v.getLinkOptions(),
	}
	cpp.SetDialectFromString(v.getDialect())
	return cpp, nil
}

func LuaCPPProfileToCPPProfile(v *LuaCPPProfile) *project.CPPProfile {
	cpp := &project.CPPProfile{
		BuildOptions: v.getBuildOptions(),
		LinkOptions:  v.getLinkOptions(),
	}
	cpp.SetDialectFromString(v.getDialect())
	return cpp
}
