package luabslib_test

import (
	"fmt"
	"testing"

	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/lua/luabslib"
	lua "github.com/yuin/gopher-lua"
)

func TestLoaderRet(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	var cppprofile *luabslib.CPPProfile
	L.PreloadModule("cppprofile", luabslib.NewCPPProfileLoader(&cppprofile))
	if err := L.DoString(`
p = require "cppprofile"
p:Dialect "toto"
p:AddBuildOptions "-Toto"
`); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	if cppprofile.FDialect != "toto" {
		fmt.Println(`cppprofile.DialectF != "toto"`)
		t.Fail()
	}
	if !functional.ListEqual(cppprofile.FBuildOptions, []string{"-Toto"}) {
		fmt.Println(`!functional.ListEqual(cppprofile.BuildOptions, []string{"-Toto"})`)
		t.Fail()
	}
}
