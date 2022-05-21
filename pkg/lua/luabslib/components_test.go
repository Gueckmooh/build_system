package luabslib_test

import (
	"fmt"
	"testing"

	"github.com/gueckmooh/bs/pkg/lua/luabslib"
	lua "github.com/yuin/gopher-lua"
)

func TestComponents1(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	luabslib.RegisterTypes(L)
	var components *luabslib.Components
	L.PreloadModule("components", luabslib.NewComponentsLoader(&components))
	if err := L.DoString(`
components = require "components"
c = components:NewComponent "pouet"
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
	for name, component := range components.FComponents {
		fmt.Println(name)
		fmt.Println(component.FType)
		for _, f := range component.FPrebuildActions {
			L.Push(f)
			L.Call(0, 0)
		}
		for _, p := range component.FProfiles {
			fmt.Println(p.CPP().FDialect)
		}
	}
}
