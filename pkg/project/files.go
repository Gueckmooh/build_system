package project

import "path/filepath"

type FilesPattern string

func (p *Project) GetRelPathForFile(file string) (string, error) {
	return filepath.Rel(p.Config.ProjectRootDirectory, file)
}
