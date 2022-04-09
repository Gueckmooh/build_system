package projectutils

import (
	"path/filepath"

	"github.com/gueckmooh/bs/pkg/lua"
	"github.com/gueckmooh/bs/pkg/project"
)

func GetProject(root string) (*project.Project, error) {
	C := lua.NewLuaContext()
	defer C.Close()
	proj, err := C.ReadProjectFile(filepath.Join(root, project.ProjectConfigFile))
	if err != nil {
		return nil, err
	}

	files, err := proj.GetComponentFiles(root)
	if err != nil {
		return nil, err
	}

	components, err := C.ReadComponentFiles(files)
	if err != nil {
		return nil, err
	}

	proj.Components = components
	proj.Config = project.GetDefaultConfig(root)

	return proj, nil
}
