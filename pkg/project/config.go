package project

import "path/filepath"

type Config struct {
	ProjectRootDirectory   string
	BuildRootDirectory     string
	BinDirectory           string
	LibDirectory           string
	ObjDirectory           string
	ExportHeadersDirectory string
}

const (
	DefaultBuildRootDirectory     = ".build"
	DefaultBinDirectory           = "bin"
	DefaultLibDirectory           = "lib"
	DefaultObjDirectory           = "obj"
	DefaultExportHeadersDirectory = "include"
)

func GetDefaultConfig(root string) *Config {
	return &Config{
		ProjectRootDirectory:   root,
		BuildRootDirectory:     DefaultBuildRootDirectory,
		BinDirectory:           DefaultBinDirectory,
		LibDirectory:           DefaultLibDirectory,
		ObjDirectory:           DefaultObjDirectory,
		ExportHeadersDirectory: DefaultExportHeadersDirectory,
	}
}

func (c *Config) GetBinDirectory() string {
	return filepath.Join(c.ProjectRootDirectory, c.BuildRootDirectory, c.BinDirectory)
}

func (c *Config) GetLibDirectory() string {
	return filepath.Join(c.ProjectRootDirectory, c.BuildRootDirectory, c.LibDirectory)
}

func (c *Config) GetExportedHeadersDirectory() string {
	return filepath.Join(c.ProjectRootDirectory, c.BuildRootDirectory,
		c.ExportHeadersDirectory)
}

func (c *Config) GetObjDirectory() string {
	return filepath.Join(c.ProjectRootDirectory, c.BuildRootDirectory, c.ObjDirectory)
}
