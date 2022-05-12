package newluabslib_test

import (
	"fmt"
	"testing"

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
p:Dialect "CPP20"
`); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println(cppprofile.DialectF)
}
