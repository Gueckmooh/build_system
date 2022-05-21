package luabslib_test

import (
	"fmt"
	"testing"

	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/lua/luabslib"
	lua "github.com/yuin/gopher-lua"
)

func TestProfile1(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	luabslib.RegisterCPPProfileType(L)
	var profile *luabslib.Profile
	L.PreloadModule("profile", luabslib.NewProfileLoader(&profile))
	if err := L.DoString(`
p = require "profile"
cpp = p:CPP()
cpp:Dialect "Toto"
cpp:AddBuildOptions "caca"
`); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	cppprofile := profile.CPP()
	if cppprofile.FDialect != "Toto" {
		fmt.Println(`cppprofile.DialectF != "Toto"`)
		t.Fail()
	}
	if !functional.ListEqual(cppprofile.FBuildOptions, []string{"caca"}) {
		fmt.Println(`!functional.ListEqual(cppprofile.BuildOptions, []string{"caca"})`)
		t.Fail()
	}
}
