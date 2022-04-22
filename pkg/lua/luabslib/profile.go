package luabslib

import (
	"fmt"

	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

var profileFunction = map[string]lua.LGFunction{
	"AddSources": newTablePusher("_sources_"),
}

func NewProfile(L *lua.LState, name string) (*lua.LTable, map[string]*lua.LTable) {
	table := L.SetFuncs(L.NewTable(), profileFunction)

	L.SetField(table, "_name_", lua.LString(name))
	L.SetField(table, "_sources_", lua.LNil)

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

	vsources := L.GetField(T, "_sources_")
	var sources []project.FilesPattern
	if vsources.Type() != lua.LTTable && vsources.Type() != lua.LTNil {
		return nil, fmt.Errorf("Error while getting component sources, unexpected type %s",
			vsources.Type().String())
	} else if vsources.Type() == lua.LTTable {
		sources = functional.ListMap(luaSTableToSTable(vsources.(*lua.LTable)),
			func(s string) project.FilesPattern { return project.FilesPattern(s) })
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
	p.Sources = sources // @todo unify this

	return p, nil
}
