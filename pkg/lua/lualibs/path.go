package lualibs

import (
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

var pathlibFunctions = map[string]lua.LGFunction{
	"Dir":  luaDir,
	"Base": luaBase,
	"Join": luaPathJoin,
}

func luaDir(L *lua.LState) int {
	path := L.ToString(1)
	dir := filepath.Dir(path)
	L.Push(lua.LString(dir))
	return 1
}

func luaBase(L *lua.LState) int {
	path := L.ToString(1)
	base := filepath.Base(path)
	L.Push(lua.LString(base))
	return 1
}

func luaPathJoin(L *lua.LState) int {
	var args []string
	n := 1
	for {
		v := L.Get(n)
		if v.Type() == lua.LTString {
			args = append(args, v.String())
		} else {
			break
		}
		n++
	}
	join := filepath.Join(args...)
	L.Push(lua.LString(join))
	return 1
}

func pathlibLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), pathlibFunctions)

	L.Push(mod)
	return 1
}
