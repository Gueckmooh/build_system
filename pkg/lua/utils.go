package lua

import (
	"fmt"
	"path/filepath"

	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

func newSetter(field string) lua.LGFunction {
	return func(L *lua.LState) int {
		self := L.ToTable(1)
		value := L.Get(2)
		L.SetField(self, field, value)
		return 0
	}
}

func newTableSetter(field string) lua.LGFunction {
	return func(L *lua.LState) int {
		self := L.ToTable(1)
		value := L.Get(2)
		var tvalue *lua.LTable
		if value.Type() == lua.LTString {
			tvalue = L.NewTable()
			tvalue.Append(value)
		} else if value.Type() == lua.LTTable {
			tvalue = value.(*lua.LTable)
		}
		L.SetField(self, field, tvalue)
		return 0
	}
}

func newTablePusher(field string) lua.LGFunction {
	return func(L *lua.LState) int {
		self := L.ToTable(1)
		value := L.Get(2)
		vtable := L.GetField(self, field)
		var ttable *lua.LTable
		if vtable.Type() != lua.LTTable {
			ttable = L.NewTable()
			L.SetField(self, field, ttable)
		} else {
			ttable = vtable.(*lua.LTable)
		}
		ttable.Append(value)
		return 0
	}
}

func newTableAppened(field string) lua.LGFunction {
	return func(L *lua.LState) int {
		self := L.ToTable(1)
		value := L.Get(2)
		vtable := L.GetField(self, field)
		var ttable *lua.LTable
		if vtable.Type() != lua.LTTable {
			ttable = L.NewTable()
			L.SetField(self, field, ttable)
		} else {
			ttable = vtable.(*lua.LTable)
		}
		if value.Type() == lua.LTTable {
			L.ForEach(value.(*lua.LTable), func(_ lua.LValue, v lua.LValue) {
				ttable.Append(v)
			})
		} else {
			ttable.Append(value)
		}
		return 0
	}
}

func luaSTableToSTable(T *lua.LTable) []string {
	var list []string
	T.ForEach(func(_, v lua.LValue) {
		if v.Type() == lua.LTString {
			list = append(list, v.String())
		}
	})
	return list
}

func luaSTableToSMap(T *lua.LTable) map[string]string {
	var list map[string]string = make(map[string]string)
	T.ForEach(func(k, v lua.LValue) {
		if v.Type() == lua.LTString && k.Type() == lua.LTString {
			// list = append(list, v.String())
			list[k.String()] = v.String()
		}
	})
	return list
}

func luaGetStringInTable(L *lua.LState, T *lua.LTable, field, desc string) (string, error) {
	vfield := L.GetField(T, field)
	if vfield.Type() != lua.LTString {
		return "", fmt.Errorf("Error while getting %s, unexpected type %s",
			desc, vfield.Type().String())
	}
	return vfield.(lua.LString).String(), nil
}

func luaGetTableInTable(L *lua.LState, T *lua.LTable, field, desc string) (*lua.LTable, error) {
	vfield := L.GetField(T, field)
	if vfield.Type() != lua.LTTable {
		return nil, fmt.Errorf("Error while getting %s, unexpected type %s",
			desc, vfield.Type().String())
	}
	return vfield.(*lua.LTable), nil
}

func luaMaybeGetStringInTable(L *lua.LState, T *lua.LTable, field, desc string) (*string, error) {
	vfield := L.GetField(T, field)
	if vfield.Type() != lua.LTString && vfield.Type() != lua.LTNil {
		return nil, fmt.Errorf("Error while getting %s, unexpected type %s",
			desc, vfield.Type().String())
	} else if vfield.Type() == lua.LTString {
		s := new(string)
		*s = vfield.String()
		return s, nil
	} else {
		return nil, nil
	}
}

func (C *LuaContext) GetProject(root string) (*project.Project, error) {
	// defer C.Close()
	proj, err := C.ReadProjectFile(filepath.Join(root, project.ProjectConfigFile))
	if err != nil {
		return nil, err
	}

	files, err := proj.GetComponentFiles(root)
	if err != nil {
		return nil, err
	}

	components, err := C.ReadComponentFiles(files)
	if err != nil {
		return nil, err
	}

	proj.Components = components
	proj.Config = project.GetDefaultConfig(root)

	return proj, nil
}
