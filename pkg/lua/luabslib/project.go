package luabslib

import (
	"fmt"

	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

//go:generate go run ./gen -i ./project.go -c Project -T ./gen/templates -P luabslib -o project_gen.go

type Project struct {
	FName            string
	FVersion         string
	FLanguages       []string
	FSources         []string
	FDefaultTarget   string
	FBaseProfile     *Profile
	FCPP             *CPPProfile
	FProfiles        map[string]*Profile
	FDefaultProfile  string
	FPlatforms       map[string]*Profile
	FDefaultPlatform string
}

func NewProject() *Project {
	baseProfile := NewProfile("Default")
	p := &Project{
		FName:            "",
		FVersion:         "",
		FLanguages:       []string{},
		FSources:         []string{},
		FDefaultTarget:   "",
		FBaseProfile:     baseProfile,
		FCPP:             baseProfile.FCPP,
		FProfiles:        make(map[string]*Profile),
		FDefaultProfile:  "",
		FPlatforms:       make(map[string]*Profile),
		FDefaultPlatform: "",
	}
	p.FProfiles["Default"] = baseProfile
	return p
}

func (p *Project) Name(name string) {
	p.FName = name
}

func (p *Project) Version(version string) {
	p.FVersion = version
}

func (p *Project) Languages(lang ...string) {
	p.FLanguages = append(p.FLanguages, lang...)
}

func (p *Project) AddSources(src string) {
	p.FSources = append(p.FSources, src)
}

func (p *Project) DefaultTarget(name string) {
	p.FDefaultTarget = name
}

func (p *Project) Profile(name string) *Profile {
	if v, ok := p.FProfiles[name]; ok {
		return v
	}
	pp := NewProfile(name)
	p.FProfiles[name] = pp
	return pp
}

func (p *Project) CPP() *CPPProfile {
	return p.FCPP
}

func (p *Project) DefaultProfile(name string) {
	p.FDefaultProfile = name
}

func (p *Project) Platforms(names ...string) {
	for _, name := range names {
		plat := NewProfile(name)
		p.FPlatforms[name] = plat
	}
}

func (p *Project) Platform(name string) (*Profile, error) {
	if v, ok := p.FPlatforms[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("Platform %s not found", name)
}

func (p *Project) DefaultPlatform(name string) {
	p.FDefaultPlatform = name
}

func NewProjectLoader(ret **Project) lua.LGFunction {
	return __NewProjectLoader(ret)
}

func ConvertLuaProjectToProject(proj *Project) *project.Project {
	var langIDs []project.LanguageID
	for _, lang := range proj.FLanguages {
		langIDs = append(langIDs, project.LanguageIDFromString(lang))
	}
	profiles := make(map[string]*project.Profile)
	platforms := make(map[string]*project.Profile)
	for name, profile := range proj.FProfiles {
		profiles[name] = ConvertLuaProfileToProfile(profile)
	}
	for name, profile := range proj.FPlatforms {
		platforms[name] = ConvertLuaProfileToProfile(profile)
	}
	pproj := &project.Project{
		Name:      proj.FName,
		Version:   proj.FVersion,
		Languages: langIDs,
		Sources: functional.ListMap(proj.FSources,
			func(s string) project.DirectoryPattern { return project.DirectoryPattern(s) }),
		DefaultTarget:   proj.FDefaultTarget,
		Profiles:        profiles,
		BaseProfile:     ConvertLuaProfileToProfile(proj.FBaseProfile),
		DefaultProfile:  proj.FDefaultProfile,
		Platforms:       platforms,
		DefaultPlatform: proj.FDefaultPlatform,
	}
	return pproj
}
