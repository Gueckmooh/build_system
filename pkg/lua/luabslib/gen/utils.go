package main

import (
	"bytes"
	"fmt"
	"path/filepath"
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

func MustExecuteTemplateFile(name string, funcs template.FuncMap, data any) string {
	fmt.Println(filepath.Join(templateDir, name))
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
