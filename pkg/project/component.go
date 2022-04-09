package project

type ComponentType int8

const (
	TypeExecutable ComponentType = iota
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
	}
	return TypeUnknown
}
