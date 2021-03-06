package fsutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gueckmooh/bs/pkg/functional"
	"github.com/gueckmooh/bs/pkg/globbing"
)

func IsFileMatched(patterns []*globbing.Pattern, filename string) bool {
	return functional.ListAnyOf(patterns,
		func(p *globbing.Pattern) bool { return p.Match(filename) })
}

func IsFileRepMatched(pattern *globbing.PatternReplace, filename string) bool {
	return pattern.Match(filename)
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

func GetMatchingRepFiles(pattern *globbing.PatternReplace, root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relpath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if !info.IsDir() && IsFileRepMatched(pattern, relpath) {
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

func RelAll(basepath string, files []string) ([]string, error) {
	var rels []string
	for _, file := range files {
		relfile, err := filepath.Rel(basepath, file)
		if err != nil {
			return nil, err
		}
		rels = append(rels, relfile)
	}
	return rels, nil
}

func FindFileUpstream(filename string, root string) (string, error) {
	if _, err := os.Stat(filepath.Join(root, filename)); os.IsNotExist(err) {
		if root != "/" {
			return FindFileUpstream(filename, filepath.Dir(root))
		} else {
			return "", fmt.Errorf("Could not find file '%s'", filename)
		}
	} else {
		return filepath.Join(root, filename), nil
	}
}

func CopyFile(from, to string) error {
	stat, err := os.Stat(from)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(to, data, stat.Mode().Perm())
	if err != nil {
		return err
	}
	return nil
}
