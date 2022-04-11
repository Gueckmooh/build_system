package project

import "path/filepath"

type Config struct {
	ProjectRootDirectory string
	BuildRootDirectory   string
	BinDirectory         string
	LibDirectory         string
	ObjDirectory         string
}

const (
	DefaultBuildRootDirectory = ".build"
	DefaultBinDirectory       = "bin"
	DefaultLibDirectory       = "lib"
	DefaultObjDirectory       = "obj"
)

func GetDefaultConfig(root string) *Config {
	return &Config{
		ProjectRootDirectory: root,
		BuildRootDirectory:   DefaultBuildRootDirectory,
		BinDirectory:         DefaultBinDirectory,
		LibDirectory:         DefaultLibDirectory,
		ObjDirectory:         DefaultObjDirectory,
	}
}

func (c *Config) GetBinDirectory() string {
	return filepath.Join(c.ProjectRootDirectory, c.BuildRootDirectory, c.BinDirectory)
}

func (c *Config) GetLibDirectory() string {
	return filepath.Join(c.ProjectRootDirectory, c.BuildRootDirectory, c.LibDirectory)
}

func (c *Config) GetObjDirectory() string {
	return filepath.Join(c.ProjectRootDirectory, c.BuildRootDirectory, c.ObjDirectory)
}
