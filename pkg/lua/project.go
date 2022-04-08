package lua

import (
	"fmt"
	"io/ioutil"

	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

var projectFunctions = map[string]lua.LGFunction{
	"Name":       newSetter("_name_"),
	"Version":    newSetter("_version_"),
	"Languages":  newTableSetter("_languages_"),
	"AddSources": newTablePusher("_sources_"),
}

func ProjectLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), projectFunctions)

	L.SetField(mod, "_name_", lua.LNil)
	L.SetField(mod, "_version_", lua.LNil)
	L.SetField(mod, "_languages_", lua.LNil)
	L.SetField(mod, "_sources_", lua.LNil)

	L.Push(mod)
	return 1
}

func ReadProjectFromLuaState(L *lua.LState) (*project.Project, error) {
	fmt.Printf("L.GetGlobal(\"project\"): %v\n", L.GetGlobal("project"))

	vproj := L.GetGlobal("project")
	if vproj.Type() != lua.LTTable {
		return nil, fmt.Errorf("Error while reading project, unexpected type %s",
			vproj.Type().String())
	}

	tproj := vproj.(*lua.LTable)

	vname := L.GetField(tproj, "_name_")
	if vname.Type() != lua.LTString {
		return nil, fmt.Errorf("Error while getting project name, unexpected type %s",
			vname.Type().String())
	}
	name := vname.(lua.LString).String()

	vversion := L.GetField(tproj, "_version_")
	if vversion.Type() != lua.LTString {
		return nil, fmt.Errorf("Error while getting project version, unexpected type %s",
			vversion.Type().String())
	}
	version := vversion.(lua.LString).String()

	vlanguages := L.GetField(tproj, "_languages_")
	if vlanguages.Type() != lua.LTTable {
		return nil, fmt.Errorf("Error while getting project languages, unexpected type %s",
			vlanguages.Type().String())
	}
	languages := functional.ListMap(luaSTableToSTable(vlanguages.(*lua.LTable)),
		project.LanguageIDFromString)

	vsources := L.GetField(tproj, "_sources_")
	if vsources.Type() != lua.LTTable {
		return nil, fmt.Errorf("Error while getting project sources, unexpected type %s",
			vsources.Type().String())
	}
	sources := functional.ListMap(luaSTableToSTable(vsources.(*lua.LTable)),
		func(s string) project.DirectoryPattern { return project.DirectoryPattern(s) })

	proj := &project.Project{
		Name:      name,
		Version:   version,
		Languages: languages,
		Sources:   sources,
	}

	return proj, nil
}

func (C *LuaContext) ReadProjectFile(filename string) (*project.Project, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error while reading file '%s':\n\t%s",
			filename, err.Error())
	}

	fmt.Println(string(data))

	if err := C.L.DoFile(filename); err != nil {
		return nil, fmt.Errorf("Error while executing file '%s':\n\t%s",
			filename, err.Error())
	}

	return ReadProjectFromLuaState(C.L)
}