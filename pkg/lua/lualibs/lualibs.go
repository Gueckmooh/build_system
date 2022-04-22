package lualibs

import lua "github.com/yuin/gopher-lua"

func LoadLibs(L *lua.LState, version string) {
	switch version {
	case "0.0.1":
		L.PreloadModule("fs", fslibLoader)
		L.PreloadModule("path", pathlibLoader)
	}
}
