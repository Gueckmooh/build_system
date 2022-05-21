package newluabslib_test

import (
	"fmt"
	"testing"

	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/lua/newluabslib"
	lua "github.com/yuin/gopher-lua"
)

func TestProfile1(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	newluabslib.RegisterCPPProfileType(L)
	var profile *newluabslib.Profile
	L.PreloadModule("profile", newluabslib.NewProfileLoader(&profile))
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
	if cppprofile.DialectF != "Toto" {
		fmt.Println(`cppprofile.DialectF != "Toto"`)
		t.Fail()
	}
	if !functional.ListEqual(cppprofile.BuildOptions, []string{"caca"}) {
		fmt.Println(`!functional.ListEqual(cppprofile.BuildOptions, []string{"caca"})`)
		t.Fail()
	}
}
