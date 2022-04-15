package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	alist "github.com/gueckmooh/bs/pkg/adjacency_list"
	"github.com/gueckmooh/bs/pkg/ccpp"
	"github.com/gueckmooh/bs/pkg/compiler"
	"github.com/gueckmooh/bs/pkg/fsutil"
	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/globbing"
	"github.com/gueckmooh/bs/pkg/project"
)

type BuildKind int8

const (
	fileSourceKind int8 = iota
	fileLinkedKind
	fileObjectKind
)

type FileDesc struct {
	name             string
	needsToBeRebuilt bool
	kind             int8
}

func newFileDesc(name string, kind int8) *FileDesc {
	return &FileDesc{
		name:             name,
		needsToBeRebuilt: false,
		kind:             kind,
	}
}

type Builder struct {
	Project          *project.Project
	buildUpstream    bool
	componentToBuild string
	component        *project.Component
	filesGraph       *alist.Graph[FileDesc, alist.AttributeNone]
	targetVertex     alist.VertexDescriptor
	filesVertices    map[string]alist.VertexDescriptor
	alwaysBuild      bool
	profile          string
	platform         string
}

func NewBuilder(p *project.Project, ctb string, opts ...BuildOption) (*Builder, error) {
	component, err := p.GetComponent(ctb)
	if err != nil {
		return nil, err
	}

	builder := &Builder{
		Project:          p,
		buildUpstream:    false,
		componentToBuild: ctb,
		component:        component,
		filesGraph:       alist.NewGraph[FileDesc, alist.AttributeNone](alist.DirectedGraph),
		filesVertices:    make(map[string]alist.VertexDescriptor),
		alwaysBuild:      false,
		profile:          "Default",
	}
	for _, opt := range opts {
		opt(builder)
	}

	return builder, nil
}

func (B *Builder) getOrCreateFileVertex(file string, kind int8) alist.VertexDescriptor {
	if v, ok := B.filesVertices[file]; ok {
		return v
	} else {
		v := B.filesGraph.AddVertex(newFileDesc(file, kind))
		B.filesVertices[file] = v
		return v
	}
}

func (B *Builder) getLinkOptionsForComponent() ([]compiler.CompilerOption, error) {
	var opts []compiler.CompilerOption
	if len(B.component.Requires) > 0 {
		opts = append(opts, compiler.WithLibraryDirectory(B.Project.Config.GetLibDirectory(true)))
	}
	for _, d := range B.component.Requires {
		dep, err := B.Project.GetComponent(d)
		if err != nil {
			return nil, err
		}
		if dep.Type != project.TypeHeaders {
			opts = append(opts, compiler.WithLibrary(dep.Name))
		}
	}
	return opts, nil
}

func (B *Builder) getIncludesOptionsForComponent() ([]compiler.CompilerOption, error) {
	var opts []compiler.CompilerOption
	includeBase := B.Project.Config.GetExportedHeadersDirectory(true)
	if B.component.Type == project.TypeLibrary {
		opts = append(opts, compiler.WithIncludeDirectory(filepath.Join(includeBase, B.component.Name)))
	}
	for _, d := range B.component.Requires {
		dep, err := B.Project.GetHeaderDirForComponent(d)
		if err != nil {
			return nil, err
		}
		opts = append(opts, compiler.WithIncludeDirectory(dep))
	}

	return opts, nil
}

func (B *Builder) getCompilerOptionsForComponent() ([]compiler.CompilerOption, error) {
	var opts []compiler.CompilerOption
	if o, err := B.getIncludesOptionsForComponent(); err != nil {
		return nil, err
	} else {
		opts = append(opts, o...)
	}
	if o, err := B.getLinkOptionsForComponent(); err != nil {
		return nil, err
	} else {
		opts = append(opts, o...)
	}

	projectProfile, err := B.Project.ComputeProfile(B.profile)
	if err != nil {
		return nil, err
	}
	componentProfile := B.component.ComputeProfile(B.profile)

	projectPlatform, err := B.Project.ComputePlatform(B.platform)
	if err != nil {
		return nil, err
	}
	componentPlatform := B.component.ComputePlatform(B.platform)
	platform := projectPlatform.Merge(componentPlatform)

	profile := projectProfile.Merge(componentProfile).Merge(platform)

	opts = append(opts, compiler.WithCPPDIalect(profile.GetCPPProfile().Dialect))
	for _, v := range profile.GetCPPProfile().BuildOptions {
		opts = append(opts, compiler.WithBuildOption(v))
	}
	for _, v := range profile.GetCPPProfile().LinkOptions {
		// fmt.Println(v)
		opts = append(opts, compiler.WithLinkOption(v))
	}
	return opts, nil
}

func (B *Builder) isBuildableNode(v alist.VertexDescriptor) bool {
	attr := B.filesGraph.GetVertexAttribute(v)
	if attr != nil && (attr.kind == fileLinkedKind || attr.kind == fileObjectKind) {
		return true
	}
	return false
}

func (B *Builder) computeWhatNeedsToBeRebuilt() (bool, error) {
	var checkNode func(alist.VertexDescriptor) error
	checkNode = func(v alist.VertexDescriptor) error {
		oe, err := B.filesGraph.OutEdges(v)
		if err != nil {
			return err
		}
		for _, ed := range oe {
			target, _ := B.filesGraph.Target(ed)
			checkNode(target)
		}
		if B.filesGraph.IsLeef(v) {
			return nil
		}

		if B.alwaysBuild && B.isBuildableNode(v) {
			B.filesGraph.GetVertexAttribute(v).needsToBeRebuilt = true
			return nil
		}

		stat, err := os.Stat(B.filesGraph.GetVertexAttribute(v).name)
		if os.IsNotExist(err) {
			B.filesGraph.GetVertexAttribute(v).needsToBeRebuilt = true
			return nil
		}

		for _, ed := range oe {
			target, _ := B.filesGraph.Target(ed)
			if B.filesGraph.GetVertexAttribute(target).needsToBeRebuilt {
				B.filesGraph.GetVertexAttribute(v).needsToBeRebuilt = true
				return nil
			}
			statTarget, _ := os.Stat(B.filesGraph.GetVertexAttribute(target).name) // @todo handle error
			if stat.ModTime().Before(statTarget.ModTime()) {
				B.filesGraph.GetVertexAttribute(v).needsToBeRebuilt = true
				return nil
			}
		}
		return nil
	}
	err := checkNode(B.targetVertex)
	if err != nil {
		return false, err
	}
	return B.filesGraph.GetVertexAttribute(B.targetVertex).needsToBeRebuilt, nil
}

func (B *Builder) computeFilesDependencies() error {
	sourceMatchers := functional.ListMap(B.component.GetSourcesForProfile(B.profile),
		func(s project.FilesPattern) *globbing.Pattern {
			return globbing.NewPattern(string(s))
		})

	sourceFiles, err := fsutil.GetMatchingFiles(sourceMatchers, B.component.Path)
	if err != nil {
		return err
	}
	sourceFiles = ccpp.FilterCPPSourceFiles(sourceFiles)
	sourceFiles, err = fsutil.RelAll(B.Project.Config.ProjectRootDirectory, sourceFiles)
	if err != nil {
		return err
	}

	var targetDir string
	switch B.component.Type {
	case project.TypeExecutable:
		targetDir = B.Project.Config.GetBinDirectory(true)
	case project.TypeLibrary:
		targetDir = B.Project.Config.GetLibDirectory(true)
	}

	targetPath := filepath.Join(targetDir, B.component.GetTargetName())
	targetVertex := B.filesGraph.AddVertex(newFileDesc(targetPath, fileLinkedKind))
	B.targetVertex = targetVertex

	for _, file := range sourceFiles {
		err := B.computeFileDependency(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (B *Builder) computeFileDependency(sourceFile string) error {
	fileWithoutSuffix := strings.TrimSuffix(sourceFile, filepath.Ext(sourceFile))
	fileWithoutSuffix, err := filepath.Abs(fileWithoutSuffix)
	if err != nil {
		return err
	}
	fileWithoutSuffix, err = filepath.Rel(B.component.Path, fileWithoutSuffix)
	if err != nil {
		return err
	}

	targetFile := filepath.Join(B.Project.Config.GetObjDirectory(true), B.component.Name,
		fileWithoutSuffix+".o")

	compilerOpts, err := B.getCompilerOptionsForComponent()
	if err != nil {
		return err
	}
	compiler := compiler.NewCompiler(compilerOpts...)
	target, sources, err := compiler.GetFileDependencies(targetFile, sourceFile)
	if err != nil {
		return err
	}

	targetVertex := B.getOrCreateFileVertex(target, fileObjectKind)

	B.filesGraph.AddEdge(B.targetVertex, targetVertex)

	for _, file := range sources {
		fileV := B.getOrCreateFileVertex(file, fileSourceKind)
		if ccpp.IsCPPSourceFile(file) {
			B.filesGraph.AddEdge(targetVertex, fileV)
		} else {
			B.filesGraph.AddEdge(targetVertex, fileV)
		}
	}
	return nil
}

func (B *Builder) getSourceToCompile(v alist.VertexDescriptor) (alist.VertexDescriptor, error) {
	oe, err := B.filesGraph.OutEdges(v)
	if err != nil {
		return 0, err
	}
	for _, ed := range oe {
		source, err := B.filesGraph.Target(ed)
		if err != nil {
			return 0, err
		}
		if ccpp.IsCPPSourceFile(B.filesGraph.GetVertexAttribute(source).name) {
			return source, nil
		}
	}
	return 0, fmt.Errorf("Could not find a source to compile for node %s",
		B.filesGraph.GetVertexAttribute(v).name)
}

func (B *Builder) Build() error {
	g := B.filesGraph
	var comp compiler.Compiler
	compilerOptions, err := B.getCompilerOptionsForComponent()
	if err != nil {
		return err
	}
	switch B.component.Type {
	case project.TypeLibrary:
		compilerOptions = append(compilerOptions, compiler.TargetLib)
	}
	comp = compiler.NewCompiler(compilerOptions...)

	var buildNode func(alist.VertexDescriptor) error
	buildNode = func(v alist.VertexDescriptor) error {
		oe, err := g.OutEdges(v)
		if err != nil {
			return err
		}

		for _, ed := range oe {
			target, err := g.Target(ed)
			if err != nil {
				return err
			}
			err = buildNode(target)
			if err != nil {
				return err
			}
		}
		if !g.GetVertexAttribute(v).needsToBeRebuilt {
			return nil
		}

		// Make sure the directory exists
		dir := filepath.Dir(g.GetVertexAttribute(v).name)
		err = fsutil.MkdirRecIfNotExist(dir)
		if err != nil {
			return err
		}

		// Compile object files
		if len(oe) > 0 && g.GetVertexAttribute(v).kind == fileObjectKind {
			source, err := B.getSourceToCompile(v)
			if err != nil {
				return err
			}
			err = comp.CompileFile(g.GetVertexAttribute(v).name, g.GetVertexAttribute(source).name)
			if err != nil {
				return err
			}
			// Link linkable file
		} else if len(oe) > 0 && g.GetVertexAttribute(v).kind == fileLinkedKind {
			var sources []string
			for _, ed := range oe {
				source, err := g.Target(ed)
				if err != nil {
					return nil
				}
				sources = append(sources, g.GetVertexAttribute(source).name)
			}
			err = comp.LinkFiles(g.GetVertexAttribute(v).name, sources...)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return buildNode(B.targetVertex)
}

func (B *Builder) tryBuildComponent() (bool, error) {
	headersExported, err := B.exportHeaders()
	if err != nil {
		return false, err
	}

	if err := B.computeFilesDependencies(); err != nil {
		return false, err
	}

	needBuild, err := B.computeWhatNeedsToBeRebuilt()
	if err != nil {
		return false, err
	}

	vertexWritterOption := alist.WithVertexLabelWritter[FileDesc, alist.AttributeNone](
		func(s *FileDesc) string {
			color := "black"
			if s.needsToBeRebuilt {
				color = "red"
			}
			return fmt.Sprintf(`[label="%s",color="%s"]`, s.name, color)
		})
	ioutil.WriteFile("/tmp/graphviz.dot", []byte(B.filesGraph.DumpGraphviz(vertexWritterOption)), 0o600)

	if needBuild {
		err = B.Build()
		if err != nil {
			return false, err
		}
	}

	return headersExported || needBuild, nil
}

func (B *Builder) BuildComponent() error {
	if B.component.Type == project.TypeUnknown {
		return fmt.Errorf("Unable to build component with unknown type %s", B.componentToBuild)
	}

	fmt.Printf("--------------- Building component '%s'...\n", B.componentToBuild)
	done, err := B.tryBuildComponent()
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
