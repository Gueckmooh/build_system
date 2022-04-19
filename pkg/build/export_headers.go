package build

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gueckmooh/bs/pkg/common/colors"
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
	return filepath.Join(B.Project.Config.GetExportedHeadersDirectory(true), B.component.Name)
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

func (B *Builder) getFilesToCopyOrRemove(copies map[string]string) (map[string]string, []string, error) {
	toCopy := make(map[string]string)
	toKeep := make(map[string]bool)
	var toRemove []string
	for f, t := range copies {
		from, to, err := B.getExportedHeaderPaths(f, t)
		if err != nil {
			return nil, nil, err
		}
		toKeep[to] = true
		_, err = os.Stat(to)
		if os.IsNotExist(err) || B.alwaysBuild {
			toCopy[from] = to
		}
	}

	if _, err := os.Stat(B.getComponentHeaderExportsDir()); !os.IsNotExist(err) {
		baseToPath := B.getComponentHeaderExportsDir()
		if err != nil {
			return nil, nil, err
		}
		err = filepath.Walk(baseToPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if _, keep := toKeep[path]; !keep {
				toRemove = append(toRemove, path)
			}
			return nil
		})
		if err != nil {
			return nil, nil, err
		}
	}
	return toCopy, toRemove, nil
}

func (B *Builder) getExportedHeaderPaths(from, to string) (string, string, error) {
	baseFromPath, err := filepath.Rel(B.Project.Config.ProjectRootDirectory, B.component.Path)
	if err != nil {
		return "", "", err
	}
	baseToPath := B.getComponentHeaderExportsDir()
	realFrom := filepath.Join(baseFromPath, from)
	realTo := filepath.Join(baseToPath, to)
	return realFrom, realTo, nil
}

func (B *Builder) doCopyFiles(copies map[string]string) error {
	for from, to := range copies {
		err := fsutil.MkdirRecIfNotExist(filepath.Dir(to))
		if err != nil {
			return err
		}
		tramp, err := makeTrampolineContent(from, to)
		if err != nil {
			return err
		}
		fmt.Printf("Writing %s%s%s\n", colors.StyleBold, to, colors.StyleReset)
		err = ioutil.WriteFile(to, []byte(tramp), 0o600)
		if err != nil {
			return err
		}
	}
	return nil
}

func (B *Builder) doRemoveFiles(removes []string) error {
	for _, file := range removes {
		fmt.Printf("Removing %s\n", file)
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (B *Builder) exportHeaders() (bool, error) {
	if B.component.ExportedHeaders == nil {
		return false, nil
	}
	allCopies := make(map[string]string)
	for k, v := range B.component.ExportedHeaders {
		p := globbing.NewPatternReplace(k, v)
		err := p.Compile()
		if err != nil {
			return false, err
		}
		files, err := fsutil.GetMatchingRepFiles(p, B.component.Path)
		if err != nil {
			return false, err
		}
		copies, err := B.getCopiesForExportHeaders(p, B.component.Path, files)
		if err != nil {
			return false, err
		}
		for k, v := range copies {
			allCopies[k] = v
		}
	}
	copies, removes, err := B.getFilesToCopyOrRemove(allCopies)
	if err != nil {
		return false, err
	}
	if len(copies) > 0 || len(removes) > 0 {
		fmt.Printf("%sExporting headers...%s\n",
			colors.ColorGray, colors.ColorReset)
		if len(copies) > 0 {
			err := B.doCopyFiles(copies)
			if err != nil {
				return true, err
			}
		}
		if len(removes) > 0 {
			err := B.doRemoveFiles(removes)
			if err != nil {
				return true, err
			}
		}
		return true, nil
	}
	return false, nil
}
