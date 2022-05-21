package newluabslib

import (
	"path/filepath"

	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

//go:generate go run ./gen -i ./component.go -c Component -T ./gen/templates -P newluabslib -o component_gen.go

var CurrentComponentFile string

type Component struct {
	FName             string
	FType             string
	FLanguages        []string
	FSources          []string
	FExportedHeaders  map[string]string
	FRequires         []string
	FProfiles         map[string]*Profile
	FBaseProfile      *Profile
	FCPP              *CPPProfile
	FPlatforms        map[string]*Profile
	FPrebuildActions  []*lua.LFunction
	FPostbuildActions []*lua.LFunction
	FComponentPath    string
}

func NewComponent(name string) *Component {
	baseProfile := NewProfile("Default")
	c := &Component{
		FName:             name,
		FType:             "",
		FLanguages:        []string{},
		FSources:          []string{},
		FExportedHeaders:  make(map[string]string),
		FRequires:         []string{},
		FProfiles:         make(map[string]*Profile),
		FBaseProfile:      baseProfile,
		FCPP:              baseProfile.FCPP,
		FPlatforms:        make(map[string]*Profile),
		FPrebuildActions:  []*lua.LFunction{},
		FPostbuildActions: []*lua.LFunction{},
		FComponentPath:    filepath.Dir(CurrentComponentFile),
	}
	c.FProfiles["Default"] = baseProfile
	return c
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

func (c *Component) ExportedHeaders(headers map[string]string) {
	for name, value := range headers {
		c.FExportedHeaders[name] = value
	}
}

func (c *Component) Requires(req ...string) {
	c.FRequires = append(c.FRequires, req...)
}

func (c *Component) Profile(name string) *Profile {
	if v, ok := c.FProfiles[name]; ok {
		return v
	}
	p := NewProfile(name)
	c.FProfiles[name] = p
	return p
}

func (c *Component) CPP() *CPPProfile {
	return c.FCPP
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

func (c *Component) AddPostbuildAction(act *lua.LFunction) {
	c.FPostbuildActions = append(c.FPostbuildActions, act)
}

func NewComponentLoader(ret **Component) lua.LGFunction {
	return __NewComponentLoader(ret)
}

func ConvertLuaComponentToComponent(comp *Component) *project.Component {
	var langIDs []project.LanguageID
	for _, lang := range comp.FLanguages {
		langIDs = append(langIDs, project.LanguageIDFromString(lang))
	}
	profiles := make(map[string]*project.Profile)
	platforms := make(map[string]*project.Profile)
	for name, profile := range comp.FProfiles {
		profiles[name] = ConvertLuaProfileToProfile(profile)
	}
	for name, profile := range comp.FPlatforms {
		platforms[name] = ConvertLuaProfileToProfile(profile)
	}
	ccomp := &project.Component{
		Name:      comp.FName,
		Languages: langIDs,
		Sources: functional.ListMap(comp.FSources,
			func(s string) project.FilesPattern { return project.FilesPattern(s) }),
		Type:             project.ComponentTypeFromString(comp.FType),
		Path:             comp.FComponentPath,
		ExportedHeaders:  comp.FExportedHeaders,
		Requires:         comp.FRequires,
		Profiles:         profiles,
		BaseProfile:      ConvertLuaProfileToProfile(comp.FBaseProfile),
		Platforms:        platforms,
		PrebuildActions:  comp.FPrebuildActions,
		PostbuildActions: comp.FPostbuildActions,
	}
	return ccomp
}
