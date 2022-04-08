package lua

import (
	"fmt"
	"io/ioutil"

	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

func luaNewComponent(L *lua.LState) int {
	self := L.ToTable(1)
	name := L.ToString(2)
	vtable := L.GetField(self, "_components_")
	var ttable *lua.LTable
	if vtable.Type() != lua.LTTable {
		ttable = L.NewTable()
		L.SetField(self, "_components_", ttable)
	} else {
		ttable = vtable.(*lua.LTable)
	}
	component := NewComponent(L, name)
	ttable.Append(component)
	L.Push(component)
	return 1
}

var (
	componentsFunctions = map[string]lua.LGFunction{
		"NewComponent": luaNewComponent,
	}
	componentFunction = map[string]lua.LGFunction{
		"Type":       newSetter("_type_"),
		"Languages":  newTableSetter("_languages_"),
		"AddSources": newTablePusher("_sources_"),
	}
)

func ComponentsLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), componentsFunctions)

	L.SetField(mod, "_components_", lua.LNil)

	L.Push(mod)
	return 1
}

func NewComponent(L *lua.LState, name string) *lua.LTable {
	table := L.SetFuncs(L.NewTable(), componentFunction)

	fmt.Println("New component", name)

	L.SetField(table, "_name_", lua.LString(name))
	L.SetField(table, "_type_", lua.LNil)
	L.SetField(table, "_languages_", lua.LNil)
	L.SetField(table, "_sources_", lua.LNil)

	return table
}

func ReadComponentFromLuaTable(L *lua.LState, T *lua.LTable) (*project.Component, error) {
	vname := L.GetField(T, "_name_")
	if vname.Type() != lua.LTString {
		return nil, fmt.Errorf("Error while getting component name, unexpected type %s",
			vname.Type().String())
	}
	name := vname.(lua.LString).String()

	vtype := L.GetField(T, "_type_")
	if vtype.Type() != lua.LTString {
		return nil, fmt.Errorf("Error while getting component version, unexpected type %s",
			vtype.Type().String())
	}
	ty := project.ComponentTypeFromString(vtype.(lua.LString).String())

	vlanguages := L.GetField(T, "_languages_")
	if vlanguages.Type() != lua.LTTable {
		return nil, fmt.Errorf("Error while getting component languages, unexpected type %s",
			vlanguages.Type().String())
	}
	languages := functional.ListMap(luaSTableToSTable(vlanguages.(*lua.LTable)),
		project.LanguageIDFromString)

	vsources := L.GetField(T, "_sources_")
	if vsources.Type() != lua.LTTable {
		return nil, fmt.Errorf("Error while getting component sources, unexpected type %s",
			vsources.Type().String())
	}
	sources := functional.ListMap(luaSTableToSTable(vsources.(*lua.LTable)),
		func(s string) project.FilesPattern { return project.FilesPattern(s) })

	proj := &project.Component{
		Name:      name,
		Languages: languages,
		Sources:   sources,
		Type:      ty,
	}

	return proj, nil
}

func ReadComponentsFromLuaState(L *lua.LState) ([]*project.Component, error) {
	vcomps := L.GetGlobal("components")
	if vcomps.Type() != lua.LTTable {
		return nil, fmt.Errorf("Error while reading component, unexpected type %s",
			vcomps.Type().String())
	}

	tcomps := vcomps.(*lua.LTable)

	vcomplist := L.GetField(tcomps, "_components_")
	if vcomplist.Type() != lua.LTTable {
		return nil, fmt.Errorf("Error while reading component list, unexpected type %s",
			vcomps.Type().String())
	}
	tcomplist := vcomplist.(*lua.LTable)

	var components []*project.Component
	var iterErr error = nil
	tcomplist.ForEach(func(_, v lua.LValue) {
		if v.Type() == lua.LTTable {
			comp, err := ReadComponentFromLuaTable(L, v.(*lua.LTable))
			if err != nil {
				iterErr = err
			}
			components = append(components, comp)
		}
	})
	if iterErr != nil {
		return nil, iterErr
	}

	return components, nil
}

func (C *LuaContext) ReadComponentFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error while reading file '%s':\n\t%s",
			filename, err.Error())
	}

	fmt.Println(string(data))

	if err := C.L.DoFile(filename); err != nil {
		return fmt.Errorf("Error while executing file '%s':\n\t%s",
			filename, err.Error())
	}

	// return ReadProjectFromLuaState(C.L)
	return nil
}

func (C *LuaContext) ReadComponentFiles(filenames []string) ([]*project.Component, error) {
	err := functional.ListTryApply(filenames,
		func(s string) error {
			return C.ReadComponentFile(s)
		})
	if err != nil {
		return nil, fmt.Errorf("Error while loading components:\n\t%s", err.Error())
	}
	return ReadComponentsFromLuaState(C.L)
}