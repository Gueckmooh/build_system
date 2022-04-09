package fsutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/globbing"
)

func IsFileMatched(patterns []*globbing.Pattern, filename string) bool {
	return functional.ListAnyOf(patterns,
		func(p *globbing.Pattern) bool { return p.Match(filename) })
}

func GetMatchingFiles(patterns []*globbing.Pattern, root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relpath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if !info.IsDir() && IsFileMatched(patterns, relpath) {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error while getting files matched from '%s'\n\t%s",
			root, err.Error())
	}
	return files, nil
}

func MkdirIfNotExist(dirname string) error {
	_, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirname, 0o755)
		if err != nil {
			return err
		}
	}
	return nil
}

func MkdirRecIfNotExist(dirname string) error {
	_, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		err = MkdirRecIfNotExist(filepath.Dir(dirname))
		if err != nil {
			return err
		}
		err = os.Mkdir(dirname, 0o755)
		if err != nil {
			return err
		}
	}
	return nil
}
