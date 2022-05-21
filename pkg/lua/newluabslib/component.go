package newluabslib

import (
	lua "github.com/yuin/gopher-lua"
)

//go:generate go run ./gen -i ./component.go -c Component -T ./gen/templates -P newluabslib -o component_gen.go

type Component struct {
	FType             string
	FLanguages        []string
	FSources          []string
	FExportedHeaders  []string
	FRequires         []string
	FProfiles         map[string]*Profile
	FPlatforms        map[string]*Profile
	FPrebuildActions  []*lua.LFunction
	FPostbuildActions []*lua.LFunction
}

func NewComponent() *Component {
	return &Component{
		FType:             "",
		FLanguages:        []string{},
		FSources:          []string{},
		FExportedHeaders:  []string{},
		FRequires:         []string{},
		FProfiles:         make(map[string]*Profile),
		FPlatforms:        make(map[string]*Profile),
		FPrebuildActions:  []*lua.LFunction{},
		FPostbuildActions: []*lua.LFunction{},
	}
}

func (c *Component) Type(ty string) {
	c.FType = ty
}

func (c *Component) Languages(langs ...string) {
	c.FLanguages = append(c.FLanguages, langs...)
}

func (c *Component) AddSources(sources ...string) {
	c.FSources = append(c.FSources, sources...)
}

func (c *Component) ExportedHeaders(headers ...string) {
	c.FExportedHeaders = append(c.FExportedHeaders, headers...)
}

func (c *Component) Requires(req ...string) {
	c.FRequires = append(c.FExportedHeaders, req...)
}

func (c *Component) Profile(name string) *Profile {
	if v, ok := c.FProfiles[name]; ok {
		return v
	}
	p := NewProfile(name)
	c.FProfiles[name] = p
	return p
}

func (c *Component) Platform(name string) *Profile {
	if v, ok := c.FPlatforms[name]; ok {
		return v
	}
	p := NewProfile(name)
	c.FPlatforms[name] = p
	return p
}

func (c *Component) AddPrebuildAction(act *lua.LFunction) {
	c.FPrebuildActions = append(c.FPrebuildActions, act)
}

func (c *Component) AddPosbuildAction(act *lua.LFunction) {
	c.FPostbuildActions = append(c.FPostbuildActions, act)
}

func NewComponentLoader(ret **Component) lua.LGFunction {
	return __NewComponentLoader(ret)
}
