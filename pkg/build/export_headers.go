package build

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/gueckmooh/bs/pkg/fsutil"
	"github.com/gueckmooh/bs/pkg/globbing"
)

func (B *Builder) getCopiesForExportHeaders(p *globbing.PatternReplace, root string, files []string) (map[string]string, error) {
	copies := make(map[string]string)
	for _, file := range files {
		relFile, err := filepath.Rel(root, file)
		if err != nil {
			return nil, err
		}
		target := p.Replace(relFile)
		copies[relFile] = target
	}
	return copies, nil
}

func (B *Builder) getComponentHeaderExportsDir() string {
	return filepath.Join(B.Project.Config.GetExportedHeadersDirectory(), B.component.Name)
}

const trampolineTemplate = `// file {{.Filename}}
#pragma once

#include "{{.Trampoline}}"
`

func makeTrampolineContent(from, to string) (string, error) {
	absFrom, err := filepath.Abs(from)
	if err != nil {
		return "", err
	}
	absTo, err := filepath.Abs(to)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(filepath.Dir(absTo), absFrom)
	if err != nil {
		return "", err
	}

	t, err := template.New("trampoline").Parse(trampolineTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	t.Execute(&buf, struct {
		Filename   string
		Trampoline string
	}{
		Filename:   to,
		Trampoline: rel,
	})

	return buf.String(), nil
}

func (B *Builder) doCopyFiles(copies map[string]string) error {
	baseFromPath, err := filepath.Rel(B.Project.Config.ProjectRootDirectory, B.component.Path)
	if err != nil {
		return err
	}
	baseToPath, err := filepath.Rel(B.Project.Config.ProjectRootDirectory, B.getComponentHeaderExportsDir())
	if err != nil {
		return err
	}
	for from, to := range copies {
		realFrom := filepath.Join(baseFromPath, from)
		realTo := filepath.Join(baseToPath, to)
		err = fsutil.MkdirRecIfNotExist(filepath.Dir(realTo))
		if err != nil {
			return err
		}
		tramp, err := makeTrampolineContent(realFrom, realTo)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(realTo, []byte(tramp), 0o600)
		if err != nil {
			return err
		}
	}
	return nil
}

func (B *Builder) exportHeaders() error {
	if B.component.ExportedHeaders == nil {
		return nil
	}
	fmt.Printf("Exporting headers for component '%s'...\n", B.component.Name)
	for k, v := range B.component.ExportedHeaders {
		p := globbing.NewPatternReplace(k, v)
		err := p.Compile()
		if err != nil {
			return err
		}
		files, err := fsutil.GetMatchingRepFiles(p, B.component.Path)
		if err != nil {
			return err
		}
		copies, err := B.getCopiesForExportHeaders(p, B.component.Path, files)
		if err != nil {
			return err
		}
		err = B.doCopyFiles(copies)
		if err != nil {
			return err
		}
	}
	return nil
}
