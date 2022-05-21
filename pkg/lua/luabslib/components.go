package luabslib

import (
	"fmt"

	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

//go:generate go run ./gen -i ./components.go -c Components -T ./gen/templates -P luabslib -o components_gen.go

type Components struct {
	FComponents map[string]*Component
}

func NewComponents() *Components {
	return &Components{
		FComponents: make(map[string]*Component),
	}
}

func (c *Components) NewComponent(name string) (*Component, error) {
	if _, ok := c.FComponents[name]; ok {
		return nil, fmt.Errorf("Cannot create component named %s it already exists", name)
	}
	cc := NewComponent(name)
	c.FComponents[name] = cc
	return cc, nil
}

func NewComponentsLoader(ret **Components) lua.LGFunction {
	return __NewComponentsLoader(ret)
}

func ConvertLuaComponentsToComponents(comps *Components) []*project.Component {
	var ccomps []*project.Component
	for _, c := range comps.FComponents {
		ccomps = append(ccomps, ConvertLuaComponentToComponent(c))
	}
	return ccomps
}
