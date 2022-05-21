package lualibs

import lua "github.com/yuin/gopher-lua"

func LoadLibs(L *lua.LState) {
	L.PreloadModule("fs", fslibLoader)
	L.PreloadModule("path", pathlibLoader)
}
