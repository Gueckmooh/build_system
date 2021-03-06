package luabslib

import lua "github.com/yuin/gopher-lua"

func RegisterTypes(L *lua.LState) {
	__RegisterCPPProfileType(L)
	__RegisterComponentType(L)
	__RegisterComponentsType(L)
	__RegisterProfileType(L)
	__RegisterProjectType(L)
	__RegisterGitRepositoryType(L)
}
