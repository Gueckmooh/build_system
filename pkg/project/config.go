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

func (c *Config) GetBinDirectory(rel bool) string {
	if rel {
		return filepath.Join(c.BuildRootDirectory, c.BinDirectory)
	} else {
		return filepath.Join(c.ProjectRootDirectory, c.BuildRootDirectory, c.BinDirectory)
	}
}

func (c *Config) GetLibDirectory(rel bool) string {
	if rel {
		return filepath.Join(c.BuildRootDirectory, c.LibDirectory)
	} else {
		return filepath.Join(c.ProjectRootDirectory, c.BuildRootDirectory, c.LibDirectory)
	}
}

func (c *Config) GetBuildDirectory(rel bool) string {
	if rel {
		return c.BuildRootDirectory
	} else {
		return filepath.Join(c.ProjectRootDirectory, c.BuildRootDirectory)
	}
}

func (c *Config) GetExportedHeadersDirectory(rel bool) string {
	if rel {
		return filepath.Join(c.BuildRootDirectory, c.ExportHeadersDirectory)
	} else {
		return filepath.Join(c.ProjectRootDirectory, c.BuildRootDirectory, c.ExportHeadersDirectory)
	}
}

func (c *Config) GetObjDirectory(rel bool) string {
	if rel {
		return filepath.Join(c.BuildRootDirectory, c.ObjDirectory)
	} else {
		return filepath.Join(c.ProjectRootDirectory, c.BuildRootDirectory, c.ObjDirectory)
	}
}
