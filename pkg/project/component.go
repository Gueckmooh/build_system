package project

import "fmt"

type ComponentType int8

const (
	TypeExecutable ComponentType = iota
	TypeLibrary
	TypeUnknown
)

type Component struct {
	Name      string
	Languages []LanguageID
	Sources   []FilesPattern
	Type      ComponentType
	Path      string
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
