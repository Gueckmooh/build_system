package newluabslib

import lua "github.com/yuin/gopher-lua"

//go:generate go run ./gen -i ./components.go -c Components -T ./gen/templates -P newluabslib -o components_gen.go

type Components struct {
	FComponents map[string]*Component
}

func NewComponents() *Components {
	return &Components{
		FComponents: make(map[string]*Component),
	}
}

func (c *Components) NewComponent(name string) *Component {
	if _, ok := c.FComponents[name]; ok {
		panic("Cannot create component named " + name + " it already exists")
	}
	cc := NewComponent()
	c.FComponents[name] = cc
	return cc
}

func NewComponentsLoader(ret **Components) lua.LGFunction {
	return __NewComponentsLoader(ret)
}
