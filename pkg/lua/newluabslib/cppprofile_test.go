package newluabslib_test

import (
	"fmt"
	"testing"

	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/lua/newluabslib"
	lua "github.com/yuin/gopher-lua"
)

func TestLoaderRet(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	var cppprofile *newluabslib.CPPProfile
	L.PreloadModule("cppprofile", newluabslib.NewCPPProfileLoader(&cppprofile))
	if err := L.DoString(`
p = require "cppprofile"
p:Dialect "toto"
p:AddBuildOptions "-Toto"
`); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	if cppprofile.DialectF != "toto" {
		fmt.Println(`cppprofile.DialectF != "toto"`)
		t.Fail()
	}
	if !functional.ListEqual(cppprofile.BuildOptions, []string{"-Toto"}) {
		fmt.Println(`!functional.ListEqual(cppprofile.BuildOptions, []string{"-Toto"})`)
		t.Fail()
	}
}
