package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

func ExecuteTemplate(name string, temp string, funcs template.FuncMap, data any) (string, error) {
	t, err := template.New(name).Funcs(funcs).Parse(temp)
	if err != nil {
		return "", err
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, data)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func MustExecuteTemplate(name string, temp string, funcs template.FuncMap, data any) string {
	t, err := template.New(name).Funcs(funcs).Parse(temp)
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

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

var basicFunctions = template.FuncMap{
	"add":      func(a, b string) string { return fmt.Sprintf("%d", must(strconv.Atoi(a))+must(strconv.Atoi(b))) },
	"join":     func(l []string, s string) string { return strings.Join(l, s) },
	"new_list": func() []string { return []string{} },
	"append":   func(sl []string, s string) []string { return append(sl, s) },
}

func MustExecuteTemplateFile(name string, funcs template.FuncMap, data any) string {
	fmt.Println(filepath.Join(templateDir, name))
	for k, v := range basicFunctions {
		funcs[k] = v
	}
	t, err := template.New(name).Funcs(funcs).ParseFiles(filepath.Join(templateDir, name))
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
