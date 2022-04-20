package build

import (
	"fmt"

	"github.com/gueckmooh/bs/pkg/lua/luadump"
	lua "github.com/yuin/gopher-lua"
)

func (B *Builder) RunLuaFunction(F *lua.LFunction) error {
	var args []lua.LValue
	for _, arg := range F.Proto.Chunk.ParList.Names {
		a, err := B.getParamForName(arg)
		if err != nil {
			return fmt.Errorf("Error while running function:\n%s\n\t%s",
				luadump.DumpFunction(F), err.Error())
		}
		args = append(args, a)
	}
	B.C.L.Push(F)
	for _, v := range args {
		B.C.L.Push(v)
	}
	B.C.L.Call(len(args), 0)
	return nil
}

func (B *Builder) getParamForName(name string) (lua.LValue, error) {
	switch name {
	case "componentName":
		return lua.LString(B.component.Name), nil
	case "targetPath":
		return lua.LString(B.filesGraph.GetVertexAttribute(B.targetVertex).name), nil
	default:
		return nil, fmt.Errorf("Unknown param '%s'", name)
	}
}
