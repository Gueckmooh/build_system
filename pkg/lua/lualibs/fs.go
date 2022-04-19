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
}

func luaCopyFile(L *lua.LState) int {
	from := L.ToString(1)
	to := L.ToString(2)
	fmt.Printf("CopyFile(%s, %s)\n", from, to)
	err := fsutil.CopyFile(from, to)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s could not copy file '%s' to '%s':\n\t%s\n",
			colors.ColorRed+colors.StyleBold, colors.ColorReset, from, to, err.Error())
	}
	return 0
}

func fslibLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), fslibFunctions)

	L.Push(mod)
	return 1
}
