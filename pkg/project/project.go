package project

import (
	"path/filepath"

	"github.com/gueckmooh/bs/pkg/fsutil"
	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/globbing"
)

const (
	ProjectConfigFile   = "bs_project.lua"
	ComponentConfigFile = "bs_component.lua"
)

type Project struct {
	Name       string
	Version    string
	Languages  []LanguageID
	Sources    []DirectoryPattern
	Components []*Component
	Config     *Config
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
