package lua

import (
	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

var cppprofileFunction = map[string]lua.LGFunction{
	"Dialect":         newSetter("_dialect_"),
	"AddBuildOptions": newTableAppened("_build_options_"),
}

func NewCPPProfile(L *lua.LState) *lua.LTable {
	table := L.SetFuncs(L.NewTable(), cppprofileFunction)

	L.SetField(table, "_dialect_", lua.LNil)
	L.SetField(table, "_build_options_", L.NewTable())

	return table
}

func ReadCPPProfileFromLuaTable(L *lua.LState, T *lua.LTable) (*project.CPPProfile, error) {
	maybeDialect, err := luaMaybeGetStringInTable(L, T, "_dialect_", "CPP dialect")
	if err != nil {
		return nil, err
	}

	vbuildOptions, err := luaGetTableInTable(L, T, "_build_options_", "CPP build options")
	if err != nil {
		return nil, err
	}

	p := project.NewCPPProfile()
	if maybeDialect != nil {
		p.SetDialectFromString(*maybeDialect)
	}
	p.BuildOptions = luaSTableToSTable(vbuildOptions)

	return p, nil
}
