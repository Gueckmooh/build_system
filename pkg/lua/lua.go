package lua

import (
	"fmt"
	"os"

	"github.com/gueckmooh/bs/pkg/lua/lualibs"
	lua "github.com/yuin/gopher-lua"
)

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

func luaSetBSVersion(L *lua.LState) int {
	version := L.ToString(1)
	if version != "0.0.1" {
		fmt.Fprintf(os.Stderr, "Unknown version '%s'\n", version)
		L.Panic(L)
	}
	LoadLuaBSLib(L)
	return 0
}

func LoadLuaBSLib(L *lua.LState) {
	L.PreloadModule("project", ProjectLoader)
	L.PreloadModule("components", ComponentsLoader)
	lualibs.LoadLibs(L)
}

func InitializeLuaState(L *lua.LState) {
	L.SetGlobal("version", L.NewFunction(luaSetBSVersion))
}
