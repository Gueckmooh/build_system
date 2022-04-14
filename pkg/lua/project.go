package lua

import (
	"fmt"

	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/project"
	lua "github.com/yuin/gopher-lua"
)

var projectFunctions = map[string]lua.LGFunction{
	"Name":           newSetter("_name_"),
	"Version":        newSetter("_version_"),
	"Languages":      newTableSetter("_languages_"),
	"AddSources":     newTablePusher("_sources_"),
	"DefaultTarget":  newSetter("_default_target_"),
	"Profile":        luaGetOrCreatetProfile,
	"DefaultProfile": newSetter("_default_profile_"),
}

func luaGetOrCreatetProfile(L *lua.LState) int {
	self := L.ToTable(1)
	name := L.ToString(2)
	vprofiles := L.GetField(self, "_profiles_")
	if vprofiles.Type() != lua.LTTable {
		// @todo warning
		return 0
	}
	profiles := vprofiles.(*lua.LTable)
	vprof := L.GetField(profiles, name)
	if vprof.Type() == lua.LTNil { // Create profile
		prof, _ := NewProfile(L, name)
		L.SetField(profiles, name, prof)
		L.Push(prof)
		return 1
	} else {
		prof := vprof.(*lua.LTable)
		L.Push(prof)
		return 1
	}
}

func ProjectLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), projectFunctions)

	L.SetField(mod, "_name_", lua.LNil)
	L.SetField(mod, "_version_", lua.LNil)
	L.SetField(mod, "_languages_", lua.LNil)
	L.SetField(mod, "_sources_", lua.LNil)
	L.SetField(mod, "_default_target_", lua.LNil)

	profile, profileMap := NewProfile(L, "Default")
	L.SetField(mod, "_base_profile_", profile)
	for k, v := range profileMap {
		L.SetField(mod, k, v)
	}
	L.SetField(mod, "_profiles_", L.NewTable())
	L.SetField(mod, "_default_profile_", lua.LNil)

	L.Push(mod)
	return 1
}

func ReadProjectFromLuaState(L *lua.LState) (*project.Project, error) {
	vproj := L.GetGlobal("project")
	if vproj.Type() != lua.LTTable {
		return nil, fmt.Errorf("Error while reading project, unexpected type %s",
			vproj.Type().String())
	}

	tproj := vproj.(*lua.LTable)

	vname := L.GetField(tproj, "_name_")
	if vname.Type() != lua.LTString {
		return nil, fmt.Errorf("Error while getting project name, unexpected type %s",
			vname.Type().String())
	}
	name := vname.(lua.LString).String()

	vversion := L.GetField(tproj, "_version_")
	if vversion.Type() != lua.LTString {
		return nil, fmt.Errorf("Error while getting project version, unexpected type %s",
			vversion.Type().String())
	}
	version := vversion.(lua.LString).String()

	vdefaultTarget := L.GetField(tproj, "_default_target_")
	defaultTarget := ""
	if vdefaultTarget.Type() != lua.LTString && vdefaultTarget.Type() != lua.LTNil {
		return nil, fmt.Errorf("Error while getting project default target, unexpected type %s",
			vdefaultTarget.Type().String())
	} else if vdefaultTarget.Type() == lua.LTString {
		defaultTarget = vdefaultTarget.(lua.LString).String()
	}

	vlanguages := L.GetField(tproj, "_languages_")
	if vlanguages.Type() != lua.LTTable {
		return nil, fmt.Errorf("Error while getting project languages, unexpected type %s",
			vlanguages.Type().String())
	}
	languages := functional.ListMap(luaSTableToSTable(vlanguages.(*lua.LTable)),
		project.LanguageIDFromString)

	vsources := L.GetField(tproj, "_sources_")
	if vsources.Type() != lua.LTTable {
		return nil, fmt.Errorf("Error while getting project sources, unexpected type %s",
			vsources.Type().String())
	}
	sources := functional.ListMap(luaSTableToSTable(vsources.(*lua.LTable)),
		func(s string) project.DirectoryPattern { return project.DirectoryPattern(s) })

	vbaseProfile := L.GetField(tproj, "_base_profile_")
	if vbaseProfile.Type() != lua.LTTable {
		return nil, fmt.Errorf("Error while getting project default profile, unexpected type %s",
			vbaseProfile.Type().String())
	}
	baseProfile, err := ReadProfileFromLuaTable(L, vbaseProfile.(*lua.LTable))
	if err != nil {
		return nil, err
	}

	profiles := make(map[string]*project.Profile)
	profiles[baseProfile.Name] = baseProfile
	{
		vprofiles := L.GetField(tproj, "_profiles_")
		if vprofiles.Type() != lua.LTTable {
			return nil, fmt.Errorf("Error while getting project profiles, unexpected type %s",
				vprofiles.Type().String())
		}
		var subProfiles []*project.Profile
		L.ForEach(vprofiles.(*lua.LTable), func(_ lua.LValue, vt lua.LValue) {
			if vt.Type() == lua.LTTable {
				profile, err := ReadProfileFromLuaTable(L, vt.(*lua.LTable))
				if err == nil {
					subProfiles = append(subProfiles, profile)
				}
			}
		})
		for _, profile := range subProfiles {
			profiles[profile.Name] = profile
			baseProfile.AddSubProfile(profile)
		}
	}

	maybeDefaultProfile, err := luaMaybeGetStringInTable(L, tproj, "_default_profile_", "default profile")
	if err != nil {
		return nil, err
	}
	defaultProfile := ""
	if maybeDefaultProfile != nil {
		defaultProfile = *maybeDefaultProfile
	}

	proj := &project.Project{
		Name:           name,
		Version:        version,
		Languages:      languages,
		Sources:        sources,
		DefaultTarget:  defaultTarget,
		BaseProfile:    baseProfile,
		Profiles:       profiles,
		DefaultProfile: defaultProfile,
	}

	return proj, nil
}

func (C *LuaContext) ReadProjectFile(filename string) (*project.Project, error) {
	if err := C.L.DoFile(filename); err != nil {
		return nil, fmt.Errorf("Error while executing file '%s':\n\t%s",
			filename, err.Error())
	}

	return ReadProjectFromLuaState(C.L)
}
