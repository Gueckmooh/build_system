package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"text/template"
)

var templateDir string

var templateFunctions = template.FuncMap{
	"add":        func(a, b int) int { return a + b },
	"appendl":    func(l []string, v string) []string { return append(l, v) },
	"newl":       func() []string { return []string{} },
	"join":       func(l []string, s string) string { return strings.Join(l, s) },
	"getPackage": func() string { return packageName },
}

func MustExecuteTemplate(filename string, data any) string {
	t, err := template.New(filename).Funcs(templateFunctions).ParseFiles(filepath.Join(templateDir, filename))
	if err != nil {
		panic(err)
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, data)
	if err != nil {
		panic(err)
	}
	return buff.String()
}
