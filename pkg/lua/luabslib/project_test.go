package luabslib_test

import (
	"fmt"
	"testing"

	"github.com/gueckmooh/bs/pkg/lua/luabslib"
	lua "github.com/yuin/gopher-lua"
)

func TestProject1(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	luabslib.RegisterTypes(L)
	var project *luabslib.Project
	L.PreloadModule("project", luabslib.NewProjectLoader(&project))
	if err := L.DoString(`
project = require "project"

project:Name    "My Pretty Project"
project:Version "0.0.1"

project:Languages     "CPP"     -- Enables C++ compilation

project:AddSources "sources/"

project:DefaultTarget "hello_exe"

`); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
}
