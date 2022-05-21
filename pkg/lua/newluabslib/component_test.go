package newluabslib_test

import (
	"fmt"
	"testing"

	"github.com/gueckmooh/bs/pkg/lua/newluabslib"
	lua "github.com/yuin/gopher-lua"
)

func TestComponent1(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	newluabslib.RegisterCPPProfileType(L)
	newluabslib.RegisterProfileType(L)
	var component *newluabslib.Component
	L.PreloadModule("component", newluabslib.NewComponentLoader(&component))
	if err := L.DoString(`
c = require "component"
c:Type "toto"
function pouet()
print("caca")
end
c:AddPrebuildAction(pouet)
p = c:Profile "zoo"
p:CPP():Dialect "CPPPP"
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
