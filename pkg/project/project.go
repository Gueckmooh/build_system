package project

import (
	"fmt"
	"path/filepath"

	alist "github.com/gueckmooh/bs/pkg/adjacency_list"
	"github.com/gueckmooh/bs/pkg/fsutil"
	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/globbing"
)

const (
	ProjectConfigFile   = "bs_project.lua"
	ComponentConfigFile = "bs_component.lua"
)

type Project struct {
	Name            string
	Version         string
	Languages       []LanguageID
	Sources         []DirectoryPattern
	Components      []*Component
	Config          *Config
	DefaultTarget   string
	ComponentDeps   *ComponentDependencyGraph
	Profiles        map[string]*Profile
	BaseProfile     *Profile
	DefaultProfile  string
	Platforms       map[string]*Profile
	DefaultPlatform string
}

type ComponentDependencyGraph struct {
	G    *alist.Graph[Component, alist.AttributeNone]
	Vmap map[*Component]alist.VertexDescriptor
}

func (p *Project) ComputeComponentDependencies() error {
	g := alist.NewGraph[Component, alist.AttributeNone](alist.DirectedGraph)
	vmap := make(map[*Component]alist.VertexDescriptor)
	getVertex := func(c *Component) alist.VertexDescriptor {
		if v, ok := vmap[c]; ok {
			return v
		} else {
			v := g.AddVertex(c)
			vmap[c] = v
			return v
		}
	}
	for _, c := range p.Components {
		v := getVertex(c)
		for _, d := range c.Requires {
			cd, err := p.GetComponent(d)
			if err != nil {
				return err
			}
			u := getVertex(cd)
			g.AddEdge(v, u)
		}
	}

	// vertexWritterOption := alist.WithVertexLabelWritter[Component, alist.AttributeNone](func(s *Component) string {
	// 	return fmt.Sprintf(`[label="%s"]`, s.Name)
	// })
	// ioutil.WriteFile("/tmp/graphviz.dot", []byte(g.DumpGraphviz(vertexWritterOption)), 0o600)
	p.ComponentDeps = &ComponentDependencyGraph{
		G:    g,
		Vmap: vmap,
	}
	return nil
}

func (p *Project) GetComponent(componentName string) (*Component, error) {
	for _, c := range p.Components {
		if c.Name == componentName {
			return c, nil
		}
	}
	return nil, fmt.Errorf("Could not find component '%s'", componentName)
}

func (p *Project) GetComponentByPath(componentPath string) (*Component, error) {
	for _, c := range p.Components {
		if c.Path == componentPath {
			return c, nil
		}
	}
	return nil, fmt.Errorf("Could not find component in '%s'", componentPath)
}

func (p *Project) GetHeaderDirForComponent(componentName string) (string, error) {
	c, err := p.GetComponent(componentName)
	if err != nil {
		return "", err
	}
	return filepath.Join(p.Config.GetExportedHeadersDirectory(true), c.Name), nil
}

func (p *Project) GetComponentFiles(root string) ([]string, error) {
	patterns := functional.ListMap(p.Sources,
		func(p DirectoryPattern) *globbing.Pattern {
			return globbing.NewPattern(string(p))
		})
	files, err := fsutil.GetMatchingFiles(patterns, root)
	if err != nil {
		return nil, err
	}
	files = functional.ListFilter(files,
		func(p string) bool {
			return filepath.Base(p) == ComponentConfigFile
		})
	return files, nil
}

func (p *Project) ComputeProfile(name string) (*Profile, error) {
	profileToMerge, ok := p.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("Could not find profile '%s'", name)
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
	return processProfile(profileToMerge), nil
}

func (p *Project) ComputePlatform(name string) (*Profile, error) {
	if name == "" {
		return DummyProfile("Default"), nil
	}
	profileToMerge, ok := p.Platforms[name]
	if !ok {
		return nil, fmt.Errorf("Could not find profile '%s'", name)
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
	return processProfile(profileToMerge), nil
}
