package lua

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gueckmooh/bs/notpkg/luabslib"
	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/lua/lualibs"
	"github.com/gueckmooh/bs/pkg/lua/newluabslib"
	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

type LuaContext struct {
	L          *lua.LState
	Project    *newluabslib.Project
	Components *newluabslib.Components
	opened     bool
}

func NewLuaContext() *LuaContext {
	C := &LuaContext{
		L:      lua.NewState(),
		opened: true,
	}
	C.InitializeLuaState()
	return C
}

func (C *LuaContext) Close() {
	if C.opened {
		C.L.Close()
		C.opened = false
	}
}

func luaSetBSVersion(L *lua.LState) int {
	version := L.ToString(1)
	if version != "0.0.1" {
		fmt.Fprintf(os.Stderr, "Unknown version '%s'\n", version)
		L.Panic(L)
	}
	// @note: for now the version is unused but it is added in
	// prevention of future releases
	return 0
}

func (C *LuaContext) LoadLuaBSLib() {
	L := C.L
	L.PreloadModule("project", newluabslib.NewProjectLoader(&C.Project))
	L.PreloadModule("components", newluabslib.NewComponentsLoader(&C.Components))
	lualibs.LoadLibs(L)
}

func (C *LuaContext) InitializeLuaState() {
	L := C.L
	L.SetGlobal("version", L.NewFunction(luaSetBSVersion))
	C.LoadLuaBSLib()
}

func (C *LuaContext) ReadComponentFile(filename string) error {
	luabslib.CurrentComponentFile = filename
	newluabslib.CurrentComponentFile = filename
	if err := C.L.DoFile(filename); err != nil {
		luabslib.CurrentComponentFile = ""
		return fmt.Errorf("Error while executing file '%s':\n\t%s",
			filename, err.Error())
	}
	luabslib.CurrentComponentFile = ""

	return nil
}

func (C *LuaContext) ReadComponentFiles(filenames []string) ([]*project.Component, error) {
	err := functional.ListTryApply(filenames,
		func(s string) error {
			return C.ReadComponentFile(s)
		})
	if err != nil {
		return nil, fmt.Errorf("Error while loading components:\n\t%s", err.Error())
	}
	return luabslib.ReadComponentsFromLuaState(C.L)
}

func (C *LuaContext) ReadProjectFile(filename string) (*project.Project, error) {
	if err := C.L.DoFile(filename); err != nil {
		return nil, fmt.Errorf("Error while executing file '%s':\n\t%s",
			filename, err.Error())
	}

	return luabslib.ReadProjectFromLuaState(C.L)
}

func (C *LuaContext) GetProject(root string) (*project.Project, error) {
	// defer C.Close()
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
