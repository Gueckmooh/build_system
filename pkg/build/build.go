package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	adjacencylist "github.com/gueckmooh/bs/pkg/adjacency_list"
	"github.com/gueckmooh/bs/pkg/ccpp"
	"github.com/gueckmooh/bs/pkg/fsutil"
	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/gcc"
	"github.com/gueckmooh/bs/pkg/globbing"
	"github.com/gueckmooh/bs/pkg/project"
)

type BuildKind int8

const (
	buildExe BuildKind = iota
	buildLib
	buildUnknown
)

type Builder struct {
	Project          *project.Project
	buildUpstream    bool
	sourcesToBuild   *SourceDependency
	buildKind        BuildKind
	componentToBuild string
	component        *project.Component
}

func BuildExe(b *Builder) {
	b.buildKind = buildExe
}

// func BuildLib(b *Builder) {
// 	b.buildKind = buildLib
// }

func NewBuilder(p *project.Project, ctb string, opts ...BuildOption) *Builder {
	builder := &Builder{
		Project:          p,
		buildUpstream:    false,
		buildKind:        buildUnknown,
		componentToBuild: ctb,
	}
	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

func (B *Builder) prepareBuildArea() error {
	err := fsutil.MkdirIfNotExist(B.Project.Config.BuildRootDirectory)
	if err != nil {
		return err
	}

	err = fsutil.MkdirIfNotExist(filepath.Join(B.Project.Config.BuildRootDirectory,
		B.Project.Config.BinDirectory))
	if err != nil {
		return err
	}

	err = fsutil.MkdirIfNotExist(filepath.Join(B.Project.Config.BuildRootDirectory,
		B.Project.Config.ObjDirectory))
	if err != nil {
		return err
	}

	return nil
}

type BuildAction int8

func buildAction(ba BuildAction) *BuildAction { return &ba }

const (
	BACompile BuildAction = iota
	BALink
	BANone
	BAUnknown
)

type FileDesc struct {
	name             string
	needsToBeRebuilt bool
	action           BuildAction
}

func NewFileDesc(name string, ba BuildAction) *FileDesc {
	return &FileDesc{
		name:             name,
		needsToBeRebuilt: false,
		action:           ba,
	}
}

type SourceDependency struct {
	g             *adjacencylist.Graph[FileDesc, BuildAction]
	target        adjacencylist.VertexDescriptor
	filesVertices map[string]adjacencylist.VertexDescriptor
	project       *project.Project
	component     *project.Component
}

func (sd *SourceDependency) GetOrAddFile(file string, ba BuildAction) (adjacencylist.VertexDescriptor, error) {
	fp, err := sd.project.GetRelPathForFile(file)
	if err != nil {
		return 0, err
	}
	if v, ok := sd.filesVertices[fp]; ok {
		return v, nil
	} else {
		v := sd.g.AddVertex(NewFileDesc(fp, ba))
		sd.filesVertices[fp] = v
		return v, nil
	}
}

func (sd *SourceDependency) getIncludesOptionsForComponent() []gcc.GCCOption {
	var opts []gcc.GCCOption
	includeBase := sd.project.Config.GetExportedHeadersDirectory()
	if sd.component.Type == project.TypeLibrary {
		opts = append(opts, gcc.WithInclude(filepath.Join(includeBase, sd.component.Name)))
	}

	return opts
}

func (sd *SourceDependency) ProcessFile(file string) error {
	fileWithoutSuffix := strings.TrimSuffix(file, filepath.Ext(file))

	base, err := filepath.Rel(sd.component.Path, fileWithoutSuffix)
	if err != nil {
		return err
	}

	obj := filepath.Join(sd.project.Config.GetObjDirectory(), sd.component.Name, base+".o")

	compiler := gcc.NewGPP(sd.getIncludesOptionsForComponent()...)
	target, sources, err := compiler.GetBuildInfoForFile(obj, file)
	if err != nil {
		return err
	}

	objV, err := sd.GetOrAddFile(target, BACompile)
	if err != nil {
		return err
	}
	sd.g.AddEdge(sd.target, objV, buildAction(BALink))

	for _, f := range sources {
		fileV, err := sd.GetOrAddFile(f, BANone)
		if err != nil {
			return err
		}
		if ccpp.IsCPPSourceFile(f) {
			sd.g.AddEdge(objV, fileV, buildAction(BACompile))
		} else {
			sd.g.AddEdge(objV, fileV, buildAction(BANone))
		}
	}
	return nil
}

func (sd *SourceDependency) CheckWhatNeedsToBeRebuilt() bool {
	var checkNode func(adjacencylist.VertexDescriptor)
	checkNode = func(v adjacencylist.VertexDescriptor) {
		oe, _ := sd.g.OutEdges(v) // @todo handle error
		for _, ed := range oe {
			target, _ := sd.g.Target(ed)
			checkNode(target)
		}
		if sd.g.IsLeef(v) {
			return
		}

		stat, err := os.Stat(sd.g.GetVertexAttribute(v).name)
		if os.IsNotExist(err) {
			sd.g.GetVertexAttribute(v).needsToBeRebuilt = true
			return
		}

		for _, ed := range oe {
			target, _ := sd.g.Target(ed)
			if sd.g.GetVertexAttribute(target).needsToBeRebuilt {
				sd.g.GetVertexAttribute(v).needsToBeRebuilt = true
				return
			}
			statTarget, _ := os.Stat(sd.g.GetVertexAttribute(target).name) // @todo handle error
			if stat.ModTime().Before(statTarget.ModTime()) {
				sd.g.GetVertexAttribute(v).needsToBeRebuilt = true
				return
			}
		}
	}
	checkNode(sd.target)
	return sd.g.GetVertexAttribute(sd.target).needsToBeRebuilt
}

func (B *Builder) GetSourcesDependencies(proj *project.Project, component *project.Component) (*SourceDependency, error) {
	sd := &SourceDependency{
		g:             adjacencylist.NewGraph[FileDesc, BuildAction](adjacencylist.DirectedGraph),
		filesVertices: make(map[string]adjacencylist.VertexDescriptor),
		project:       proj,
		component:     component,
	}
	g := sd.g

	srcMatchers := functional.ListMap(component.Sources, func(s project.FilesPattern) *globbing.Pattern {
		return globbing.NewPattern(string(s))
	})

	files, err := fsutil.GetMatchingFiles(srcMatchers, component.Path)
	if err != nil {
		return nil, fmt.Errorf("Error while getting sources files of component %s\n\t%s",
			component.Name, err.Error())
	}
	files = ccpp.FilterCPPSourceFiles(files)

	var buildDir string
	switch component.Type {
	case project.TypeExecutable:
		buildDir = B.Project.Config.GetBinDirectory()
	case project.TypeLibrary:
		buildDir = B.Project.Config.GetLibDirectory()
	}
	targetName := filepath.Join(buildDir, component.GetTargetName())
	targetName, err = sd.project.GetRelPathForFile(targetName)
	if err != nil {
		return nil, err
	}
	target := g.AddVertex(NewFileDesc(targetName, BALink))
	sd.target = target
	for _, f := range files {
		err := sd.ProcessFile(f)
		if err != nil {
			return nil, err
		}
	}

	return sd, nil
}

func (B *Builder) Build() error {
	g := B.sourcesToBuild.g
	var compiler *gcc.GCC
	compilerOptions := B.sourcesToBuild.getIncludesOptionsForComponent()
	switch B.buildKind {
	case buildLib:
		compilerOptions = append(compilerOptions, gcc.TargetLib)
	}
	compiler = gcc.NewGPP(compilerOptions...)
	var buildNode func(adjacencylist.VertexDescriptor) error
	buildNode = func(v adjacencylist.VertexDescriptor) error {
		oe, err := g.OutEdges(v)
		if err != nil {
			return err
		}

		for _, ed := range oe {
			target, _ := g.Target(ed)
			err := buildNode(target)
			if err != nil {
				return err
			}
		}
		if g.IsLeef(v) || !g.GetVertexAttribute(v).needsToBeRebuilt {
			return nil
		}

		dir := filepath.Dir(g.GetVertexAttribute(v).name)
		err = fsutil.MkdirRecIfNotExist(dir)
		if err != nil {
			return err
		}

		if len(oe) > 0 && g.GetVertexAttribute(v).action == BACompile {
			var source adjacencylist.VertexDescriptor
			for _, ed := range oe {
				if *g.GetEdgeAttribute(ed) == BACompile {
					source, err = g.Target(ed)
					if err != nil {
						return err
					}
					break
				}
			}
			err := compiler.CompileFile(g.GetVertexAttribute(v).name, g.GetVertexAttribute(source).name)
			if err != nil {
				return err
			}
		} else if len(oe) > 0 && g.GetVertexAttribute(v).action == BALink {
			var sources []string
			for _, ed := range oe {
				source, err := g.Target(ed)
				if err != nil {
					return nil
				}
				sources = append(sources, g.GetVertexAttribute(source).name)
			}
			compiler.LinkFile(g.GetVertexAttribute(v).name, sources...)
		}
		return nil
	}
	err := buildNode(B.sourcesToBuild.target)
	return err
}

func (B *Builder) tryBuildComponent(component *project.Component) (bool, error) {
	headersHasBeenExported, err := B.exportHeaders()
	if err != nil {
		return false, err
	}

	if err := B.prepareBuildArea(); err != nil {
		return false, fmt.Errorf("Fail to prepare build area:\n\t%s", err.Error())
	}

	srcDeps, err := B.GetSourcesDependencies(B.Project, component)
	if err != nil {
		return false, err
	}

	needBuild := srcDeps.CheckWhatNeedsToBeRebuilt()
	B.sourcesToBuild = srcDeps

	vertexWritterOption := adjacencylist.WithVertexLabelWritter[FileDesc, BuildAction](func(s *FileDesc) string {
		color := "black"
		if s.needsToBeRebuilt {
			color = "red"
		}
		return fmt.Sprintf(`[label="%s",color="%s"]`, s.name, color)
	})
	ioutil.WriteFile("/tmp/graphviz.dot", []byte(srcDeps.g.DumpGraphviz(vertexWritterOption)), 0o600)

	if needBuild {
		err = B.Build()
		if err != nil {
			return false, err
		}
	}

	return needBuild || headersHasBeenExported, nil
}

func (B *Builder) BuildComponent() error {
	var component *project.Component
	if mc := functional.ListFindIf(B.Project.Components, func(c *project.Component) bool {
		return c.Name == B.componentToBuild
	}); mc != nil {
		component = *mc
	} else {
		return fmt.Errorf("Could not find component %s", B.componentToBuild)
	}
	B.component = component

	switch component.Type {
	case project.TypeExecutable:
		B.buildKind = buildExe
	case project.TypeLibrary:
		B.buildKind = buildLib
	case project.TypeUnknown:
		return fmt.Errorf("Unable to build component with unknown type %s", B.componentToBuild)
	}
	fmt.Printf("--------------- Building component '%s'...\n", B.componentToBuild)
	done, err := B.tryBuildComponent(component)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while building componnent '%s':\n\t%s\n", B.componentToBuild, err.Error())
		fmt.Printf("--------------- Failed to build component '%s'\n", B.componentToBuild)
		return fmt.Errorf("Build of component '%s' failed", B.componentToBuild)
	}
	if done {
		fmt.Printf("--------------- Build successful\n")
	} else {
		fmt.Printf("--------------- Nothing to be done for '%s'\n", B.componentToBuild)
	}
	return nil
}
