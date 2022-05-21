package luabslib_test

import (
	"fmt"
	"testing"

	"github.com/gueckmooh/bs/pkg/lua/luabslib"
	lua "github.com/yuin/gopher-lua"
)

func TestComponent1(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	luabslib.RegisterCPPProfileType(L)
	luabslib.RegisterProfileType(L)
	var component *luabslib.Component
	L.PreloadModule("component", luabslib.NewComponentLoader(&component))
	if err := L.DoString(`
c = require "component"
c:Type "toto"
function pouet()
print("caca")
end
c:AddPrebuildAction(pouet)
p = c:Profile "zoo"
p:CPP():Dialect "CPPPP"
c:ExportedHeaders {
  ["src/[DIRS]/*.hpp"] = "debug/[DIRS]/*.hpp",
}

`); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println(component.FType)
	for _, f := range component.FPrebuildActions {
		L.Push(f)
		L.Call(0, 0)
	}
	for _, p := range component.FProfiles {
		fmt.Println(p.CPP().FDialect)
	}
}
