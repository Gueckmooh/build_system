package lua

import lua "github.com/yuin/gopher-lua"

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
