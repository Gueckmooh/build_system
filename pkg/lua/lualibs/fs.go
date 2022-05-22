package lualibs

import (
	"fmt"
	"os"

	"github.com/gueckmooh/bs/pkg/common/colors"
	"github.com/gueckmooh/bs/pkg/fsutil"
	lua "github.com/yuin/gopher-lua"
)

var fslibFunctions = map[string]lua.LGFunction{
	"CopyFile": luaCopyFile,
	"Exists":   luaExists,
}

func luaCopyFile(L *lua.LState) int {
	from := L.ToString(1)
	to := L.ToString(2)
	fmt.Printf("Copying file %s%s%s to %s%s%s...\n",
		colors.StyleBold, from, colors.StyleReset,
		colors.StyleBold, to, colors.StyleReset)
	err := fsutil.CopyFile(from, to)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s could not copy file '%s' to '%s':\n\t%s\n",
			colors.ColorRed+colors.StyleBold, colors.ColorReset, from, to, err.Error())
	}
	return 0
}

func luaExists(L *lua.LState) int {
	file := L.CheckString(1)
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		L.Push(lua.LFalse)
		return 1
	}
	L.Push(lua.LTrue)
	return 1
}

func fslibLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), fslibFunctions)

	L.Push(mod)
	return 1
}
