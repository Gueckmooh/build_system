package lua

import lua "github.com/yuin/gopher-lua"

type LuaContext struct {
	L      *lua.LState
	opened bool
}

func NewLuaContext() *LuaContext {
	C := &LuaContext{
		L:      lua.NewState(),
		opened: true,
	}
	InitializeLuaState(C.L)
	return C
}

func (C *LuaContext) Close() {
	if C.opened {
		C.L.Close()
		C.opened = false
	}
}

func InitializeLuaState(L *lua.LState) {
	L.PreloadModule("project", ProjectLoader)
	L.PreloadModule("components", ComponentsLoader)
}
