package newluabslib

import lua "github.com/yuin/gopher-lua"

//go:generate go run ./gen -i ./project.go -c Project -T ./gen/templates -P newluabslib -o project_gen.go

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
	p := &Project{
		FName:            "",
		FVersion:         "",
		FLanguages:       []string{},
		FSources:         []string{},
		FDefaultTarget:   "",
		FProfiles:        make(map[string]*Profile),
		FDefaultProfile:  "",
		FPlatforms:       make(map[string]*Profile),
		FDefaultPlatform: "",
	}
	p.FBaseProfile = NewProfile("Default")
	p.FCPP = p.FBaseProfile.FCPP
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

func (p *Project) DefaultProfile(name string) {
	p.FDefaultProfile = name
}

func (p *Project) Platforms(names ...string) {
	for _, name := range names {
		plat := NewProfile(name)
		p.FPlatforms[name] = plat
	}
}

func (p *Project) Platform(name string) *Profile {
	if v, ok := p.FPlatforms[name]; ok {
		return v
	}
	panic("Platform " + name + " not found")
}

func (p *Project) DefaultPlatform(name string) {
	p.FDefaultPlatform = name
}

func NewProjectLoader(ret **Project) lua.LGFunction {
	return __NewProjectLoader(ret)
}
