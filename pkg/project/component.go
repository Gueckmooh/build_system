package project

import "fmt"

type ComponentType int8

const (
	TypeExecutable ComponentType = iota
	TypeLibrary
	TypeUnknown
)

type Component struct {
	Name            string
	Languages       []LanguageID
	Sources         []FilesPattern
	Type            ComponentType
	Path            string
	ExportedHeaders map[string]string
	Requires        []string
	Profiles        map[string]*Profile
	BaseProfile     *Profile
	Platforms       map[string]*Profile
}

func ComponentTypeFromString(compTy string) ComponentType {
	switch compTy {
	case "executable":
		return TypeExecutable
	case "library":
		return TypeLibrary
	}
	return TypeUnknown
}

func (c *Component) GetTargetName() string {
	if c.Type == TypeLibrary {
		return fmt.Sprintf("lib%s.so", c.Name)
	} else {
		return c.Name
	}
}

func (c *Component) ComputeProfile(name string) *Profile {
	profileToMerge, ok := c.Profiles[name]
	if !ok {
		profileToMerge = c.BaseProfile
	}
	var processProfile func(p *Profile) *Profile
	processProfile = func(p *Profile) *Profile {
		if p.parentProfile == nil {
			return p.Clone()
		} else {
			pp := processProfile(p.parentProfile)
			return pp.Merge(p)
		}
	}
	return processProfile(profileToMerge)
}

// @todo simplify this
func (c *Component) ComputePlatform(name string) *Profile {
	profileToMerge, ok := c.Platforms[name]
	if !ok {
		return DummyProfile("Default")
	}
	var processProfile func(p *Profile) *Profile
	processProfile = func(p *Profile) *Profile {
		if p.parentProfile == nil {
			return p.Clone()
		} else {
			pp := processProfile(p.parentProfile)
			return pp.Merge(p)
		}
	}
	return processProfile(profileToMerge)
}

func (c *Component) GetSourcesForProfile(name string) []FilesPattern {
	profileToMerge, ok := c.Profiles[name]
	if !ok {
		profileToMerge = c.BaseProfile
	}
	var processProfile func(p *Profile) *Profile
	processProfile = func(p *Profile) *Profile {
		if p.parentProfile == nil {
			return p.Clone()
		} else {
			pp := processProfile(p.parentProfile)
			return pp.Merge(p)
		}
	}
	mergedProfile := processProfile(profileToMerge)
	return append(c.Sources, mergedProfile.Sources...)
}

// @todo simplify this
func (c *Component) GetSourcesForPlatform(name string) []FilesPattern {
	profileToMerge, ok := c.Platforms[name]
	if !ok {
		return c.Sources
	}
	var processProfile func(p *Profile) *Profile
	processProfile = func(p *Profile) *Profile {
		if p.parentProfile == nil {
			return p.Clone()
		} else {
			pp := processProfile(p.parentProfile)
			return pp.Merge(p)
		}
	}
	mergedProfile := processProfile(profileToMerge)
	return append(c.Sources, mergedProfile.Sources...)
}
