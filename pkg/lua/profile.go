package lua

import (
	"fmt"

	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

var profileFunction = map[string]lua.LGFunction{}

func NewProfile(L *lua.LState, name string) (*lua.LTable, map[string]*lua.LTable) {
	table := L.SetFuncs(L.NewTable(), profileFunction)

	L.SetField(table, "_name_", lua.LString(name))

	profileMap := make(map[string]*lua.LTable)
	profileMap["CPP"] = NewCPPProfile(L)
	for k, v := range profileMap {
		L.SetField(table, k, v)
	}

	return table, profileMap
}

func ReadProfileFromLuaTable(L *lua.LState, T *lua.LTable) (*project.Profile, error) {
	name, err := luaGetStringInTable(L, T, "_name_", "profile name")
	if err != nil {
		return nil, err
	}

	p := project.NewProfile(name)

	vcppprofile := L.GetField(T, "CPP")
	if vcppprofile.Type() != lua.LTTable {
		return nil, fmt.Errorf("Error while getting CPP profile from profile %s, unexpected type %s",
			name, vcppprofile.Type().String())
	}
	cppprofile, err := ReadCPPProfileFromLuaTable(L, vcppprofile.(*lua.LTable))
	if err != nil {
		return nil, err
	}

	p.SetCPPProfile(cppprofile)

	return p, nil
}
